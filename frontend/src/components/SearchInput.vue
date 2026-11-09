<script setup lang="ts">
import { ref, watch } from 'vue';
import Search from '@primeicons/vue/search';

const props = withDefaults(
  defineProps<{
    placeholder?: string;
    debounce?: number;
  }>(),
  {
    placeholder: 'Buscar…',
    debounce: 300,
  }
);

const emit = defineEmits<{
  search: [value: string];
}>();

const query = ref('');
let timer: ReturnType<typeof setTimeout> | null = null;

watch(query, (value) => {
  if (timer) clearTimeout(timer);
  timer = setTimeout(() => emit('search', value.trim()), props.debounce);
});

function clear(): void {
  query.value = '';
  emit('search', '');
}
</script>

<template>
  <div class="flex items-center gap-2">
    <IconField class="w-full max-w-sm">
      <InputIcon>
        <Search :size="16" />
      </InputIcon>
      <InputText v-model="query" :placeholder="placeholder" class="w-full" />
    </IconField>
    <Button v-if="query" severity="secondary" label="Limpiar" @click="clear" />
  </div>
</template>
