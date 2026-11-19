<script setup lang="ts">
  import { onBeforeUnmount, ref, watch } from 'vue';
  import { Search } from '@lucide/vue';
  import { Button } from '@/components/ui/button';
  import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group';

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

  onBeforeUnmount(() => {
    if (timer) clearTimeout(timer);
  });

  function clear(): void {
    query.value = '';
    emit('search', '');
  }
</script>

<template>
  <div class="flex w-full max-w-md items-center gap-2">
    <InputGroup>
      <InputGroupAddon>
        <Search />
      </InputGroupAddon>
      <InputGroupInput v-model="query" :placeholder="placeholder" />
    </InputGroup>
    <Button v-if="query" variant="secondary" @click="clear">Limpiar</Button>
  </div>
</template>
