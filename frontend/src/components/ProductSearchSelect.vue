<script setup lang="ts">
  import { computed, onMounted, shallowRef, useTemplateRef, watch } from 'vue';
  import type {
    AutoCompleteCompleteEvent,
    AutoCompleteOptionSelectEvent,
  } from 'primevue/autocomplete';
  import Search from '@primeicons/vue/search';
  import { useProductSearch } from '@/composables/useProductSearch';
  import { formatPrice } from '@/utils/currency';
  import type { Product } from '@/types/product';
           
  const { placeholder = 'Buscar por nombre, código o código de barras…', debounce = 1000 } =
    defineProps<{
      placeholder?: string;
      /** Milliseconds to wait after the last keystroke before searching. */
      debounce?: number;
    }>();

  const emit = defineEmits<{
    select: [product: Product];
  }>();

  /** Height of one option row; the virtual scroller needs it fixed. */
  const OPTION_HEIGHT = 56;
  /** Load the next page when this close to the end of the list. */
  const LOAD_MORE_THRESHOLD = 5;

  const { suggestions, term, loading, pending, error, hasMore, search, loadMore, markStale, reset } =
    useProductSearch();

  // AutoComplete's model: the typed text, or briefly the picked product.
  const query = shallowRef<string | Product>('');
  const typedText = computed(() => (typeof query.value === 'string' ? query.value.trim() : ''));

  // `show`, `hide` and `searchTimeout` (the pending debounce) exist on the
  // component instance but aren't part of PrimeVue's typings.
  const autoComplete = useTemplateRef<{
    show: () => void;
    hide: () => void;
    searchTimeout?: ReturnType<typeof setTimeout> | null;
  }>('autoComplete');

  /** A TAB pressed while searching, waiting for the results to pick the first one. */
  const tabQueued = shallowRef(false);

  const firstSelectable = computed(() => suggestions.value.find(isSelectable));
  const searching = computed(() => pending.value || loading.value);
  const emptyMessage = computed(() => {
    if (error.value) return error.value;
    if (searching.value) return 'Buscando…';
    return `Sin coincidencias para «${term.value}»`;
  });

  const virtualScrollerOptions = {
    itemSize: OPTION_HEIGHT,
    onScrollIndexChange: ({ last }: { last: number }) => {
      if (hasMore.value && last >= suggestions.value.length - LOAD_MORE_THRESHOLD) loadMore();
    },
  };

  // Results belong to the text that produced them: once it changes they're stale.
  // A queued TAB belongs to that text too.
  watch(typedText, (text) => {
    tabQueued.value = false;
    if (text) markStale();
    else reset();
  });

  // Resolve a queued TAB once the search settles. `post` runs after
  // AutoComplete has reacted to the new suggestions, so `hide` sticks.
  watch(
    searching,
    (busy) => {
      if (busy || !tabQueued.value) return;
      tabQueued.value = false;
      const product = firstSelectable.value;
      if (product) pick(product);
      // Nothing to pick: show why (no matches, no stock, error). Needed when
      // the TAB started the search, since AutoComplete didn't and won't open.
      else autoComplete.value?.show();
    },
    { flush: 'post' }
  );

  function isSelectable(product: Product): boolean {
    return product.active && product.stock > 0;
  }

  function isDisabled(product: Product): boolean {
    return !isSelectable(product);
  }

  function stockLabel(product: Product): string {
    if (!product.active) return 'Inactivo';
    if (product.stock <= 0) return 'Sin stock';
    return `Stock: ${product.stock}`;
  }

  function select(product: Product): void {
    emit('select', product);
    query.value = '';
  }

  /** Select from the keyboard and close the overlay; focus stays in the input. */
  function pick(product: Product): void {
    select(product);
    autoComplete.value?.hide();
  }

  function onComplete(event: AutoCompleteCompleteEvent): void {
    const text = event.query.trim();
    // A late debounce for text that is no longer in the input.
    if (text !== typedText.value) return;
    // Already searched (or being searched) early by a TAB.
    if (text === term.value && (loading.value || !pending.value)) return;
    void search(text);
  }

  function onOptionSelect(event: AutoCompleteOptionSelectEvent): void {
    select(event.value as Product);
  }

  /**
   * TAB picks the first available result. It runs in the capture phase so
   * AutoComplete's own TAB handling (which would pick the highlighted option)
   * never sees the event. Focus stays in the input for the next search.
   *
   * While the search is still on its way the TAB is queued and resolved as
   * soon as the results arrive; during the debounce the search starts at once.
   */
  function onTab(event: KeyboardEvent): void {
    if (event.shiftKey) return;

    const product = firstSelectable.value;
    if (product) {
      event.preventDefault();
      event.stopPropagation();
      pick(product);
      return;
    }

    if (!typedText.value || !searching.value) return;

    event.preventDefault();
    event.stopPropagation();
    tabQueued.value = true;

    // Still debouncing: skip the wait and search now.
    if (pending.value && !loading.value) {
      const timer = autoComplete.value?.searchTimeout;
      if (timer) clearTimeout(timer);
      void search(typedText.value);
    }
  }


onMounted(() => {
  document.getElementById('mi-autocomplete')?.focus()
})
</script>

<template>
  <div @keydown.tab.capture="onTab">
    <IconField>
      <InputIcon>
        <Search :size="16" />
      </InputIcon>
      <AutoComplete
       inputId="mi-autocomplete"
        ref="autoComplete"
        v-model="query"
        :suggestions="suggestions"
        option-label="name"
        :option-disabled="isDisabled"
        :delay="debounce"
        :loading="loading"
        :placeholder="placeholder"
        :virtual-scroller-options="virtualScrollerOptions"
        scroll-height="20rem"
        fluid
        @complete="onComplete"
        @option-select="onOptionSelect"
      >
        <template #option="{ option }">
          <div class="flex w-full items-center justify-between gap-4">
            <div class="flex min-w-0 flex-col">
              <span class="truncate font-medium">{{ option.name }}</span>
              <span class="truncate font-mono text-xs text-surface-500">
                {{ option.code }}<template v-if="option.barcode"> · {{ option.barcode }}</template>
              </span>
            </div>
            <div class="flex shrink-0 flex-col items-end">
              <span class="text-sm">{{ formatPrice(option.price) }}</span>
              <span class="text-xs text-surface-500">{{ stockLabel(option) }}</span>
            </div>
          </div>
        </template>

        <template #empty>
          <span :class="{ 'text-red-500': error }">{{ emptyMessage }}</span>
        </template>

        <template #footer>
          <div
            v-if="loading && suggestions.length > 0"
            class="px-3 py-2 text-center text-xs text-surface-500"
          >
            Cargando más productos…
          </div>
        </template>
      </AutoComplete>
    </IconField>
  </div>
</template>
