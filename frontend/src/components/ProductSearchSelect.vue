<script setup lang="ts">
  import { computed, onBeforeUnmount, shallowRef, watch } from 'vue';
  import { ComboboxInput, ComboboxVirtualizer } from 'reka-ui';
  import { Search } from '@lucide/vue';
  import {
    Combobox,
    ComboboxAnchor,
    ComboboxItem,
    ComboboxList,
    ComboboxViewport,
  } from '@/components/ui/combobox';
  import { InputGroup, InputGroupAddon } from '@/components/ui/input-group';
  import { Spinner } from '@/components/ui/spinner';
  import { useProductSearch } from '@/composables/useProductSearch';
  import { formatPrice } from '@/utils/currency';
  import type { Product } from '@/types/product';

  const { placeholder = 'Buscar por nombre, código o código de barras…', debounce = 600 } =
    defineProps<{
      placeholder?: string;
      /** Milliseconds to wait after the last keystroke before searching. */
      debounce?: number;
    }>();

  const emit = defineEmits<{
    select: [product: Product];
  }>();

  /** Height of one option row; the virtualizer needs it fixed. */
  const OPTION_HEIGHT = 56;
  /** Load the next page when this close to the end of the list. */
  const LOAD_MORE_THRESHOLD = 5;

  const {
    suggestions,
    term,
    loading,
    pending,
    error,
    hasMore,
    search,
    loadMore,
    markStale,
    reset,
  } = useProductSearch();

  /** The text in the input. */
  const query = shallowRef('');
  const typedText = computed(() => query.value.trim());
  const open = shallowRef(false);

  /** A TAB pressed while searching, waiting for the results to pick the first one. */
  const tabQueued = shallowRef(false);
  let debounceTimer: ReturnType<typeof setTimeout> | undefined;

  const firstSelectable = computed(() => suggestions.value.find(isSelectable));
  const searching = computed(() => pending.value || loading.value);
  const emptyMessage = computed(() => {
    if (error.value) return error.value;
    if (searching.value) return 'Buscando…';
    return `Sin coincidencias para «${term.value}»`;
  });

  // Results belong to the text that produced them: once it changes they're
  // stale, and so is a queued TAB. The new text is searched after the debounce.
  watch(typedText, (text) => {
    tabQueued.value = false;
    clearTimeout(debounceTimer);

    if (!text) {
      reset();
      open.value = false;
      return;
    }

    markStale();
    debounceTimer = setTimeout(() => runSearch(text), debounce);
  });

  // Resolve a queued TAB once the search settles.
  watch(searching, (busy) => {
    if (busy || !tabQueued.value) return;
    tabQueued.value = false;
    const product = firstSelectable.value;
    if (product) pick(product);
    // Nothing to pick: show why (no matches, no stock, error).
    else open.value = true;
  });

  onBeforeUnmount(() => clearTimeout(debounceTimer));

  function runSearch(text: string): void {
    // A late debounce for text that is no longer in the input.
    if (text !== typedText.value) return;
    void search(text);
  }

  function isSelectable(product: Product): boolean {
    return product.active && product.stock > 0;
  }

  function stockLabel(product: Product): string {
    if (!product.active) return 'Inactivo';
    if (product.stock <= 0) return 'Sin stock';
    return `Stock: ${product.stock}`;
  }

  /** Emits the product and clears the input; focus stays in it for the next search. */
  function pick(product: Product): void {
    emit('select', product);
    query.value = '';
    open.value = false;
  }

  function onSelect(value: unknown): void {
    if (value) pick(value as Product);
  }

  function onScroll(event: Event): void {
    const el = event.target as HTMLElement;
    const remaining = el.scrollHeight - el.scrollTop - el.clientHeight;
    if (hasMore.value && remaining < LOAD_MORE_THRESHOLD * OPTION_HEIGHT) loadMore();
  }

  /**
   * TAB picks the first available result (capture phase, so the combobox never
   * sees it). While the search is still on its way the TAB is queued and
   * resolved as soon as the results arrive; during the debounce the search
   * starts at once.
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
      clearTimeout(debounceTimer);
      void search(typedText.value);
    }
  }
</script>

<template>
  <div @keydown.tab.capture="onTab">
    <!-- Always unselected: picking a product emits it and clears the input. -->
    <Combobox
      v-model:open="open"
      :model-value="null"
      ignore-filter
      :reset-search-term-on-blur="false"
      @update:model-value="onSelect"
    >
      <ComboboxAnchor class="w-full">
        <InputGroup>
          <InputGroupAddon>
            <Search />
          </InputGroupAddon>
          <ComboboxInput
            v-model="query"
            auto-focus
            data-slot="input-group-control"
            :placeholder="placeholder"
            class="h-full min-w-0 flex-1 bg-transparent px-1.5 text-sm outline-none placeholder:text-muted-foreground"
          />
          <InputGroupAddon v-if="loading" align="inline-end">
            <Spinner />
          </InputGroupAddon>
        </InputGroup>
      </ComboboxAnchor>

      <ComboboxList v-if="typedText" align="start" class="max-h-none">
        <div
          v-if="suggestions.length === 0"
          class="px-3 py-6 text-center text-sm"
          :class="error ? 'text-destructive' : 'text-muted-foreground'"
        >
          {{ emptyMessage }}
        </div>

        <ComboboxViewport v-else class="max-h-80" @scroll="onScroll">
          <ComboboxVirtualizer
            v-slot="{ option }"
            :options="suggestions"
            :estimate-size="OPTION_HEIGHT"
            :text-content="(product: Product) => product.name"
          >
            <ComboboxItem
              :value="option"
              :disabled="!isSelectable(option)"
              class="w-full pr-2"
              :style="{ height: `${OPTION_HEIGHT}px` }"
            >
              <div class="flex w-full items-center justify-between gap-4">
                <div class="flex min-w-0 flex-col">
                  <span class="truncate font-medium">{{ option.name }}</span>
                  <span class="truncate font-mono text-xs text-muted-foreground">
                    {{ option.code
                    }}<template v-if="option.barcode"> · {{ option.barcode }}</template>
                  </span>
                </div>
                <div class="flex shrink-0 flex-col items-end">
                  <span class="text-sm">{{ formatPrice(option.price) }}</span>
                  <span class="text-xs text-muted-foreground">{{ stockLabel(option) }}</span>
                </div>
              </div>
            </ComboboxItem>
          </ComboboxVirtualizer>
        </ComboboxViewport>

        <div
          v-if="loading && suggestions.length > 0"
          class="border-t px-3 py-2 text-center text-xs text-muted-foreground"
        >
          Cargando más productos…
        </div>
      </ComboboxList>
    </Combobox>
  </div>
</template>
