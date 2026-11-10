<script setup lang="ts">
  import { Plus } from '@primeicons/vue';
  import { ref } from 'vue';
  import { useCategories } from '@/composables/useCategories';
  import { useToast } from 'primevue/usetoast';

  const { create, error } = useCategories();
  const toast = useToast();

  interface Props {
    visible?: boolean;
  }
  const { visible = false } = defineProps<Props>();

  const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void;
    (e: 'created'): void;
  }>();

  const open = ref(visible);

  const form = ref({
    name: '',
  });

  async function handleCreate() {
    const success = await create(form.value.name);

    if (!success) {
      return;
    }

    form.value.name = '';
    open.value = false;
    emit('created');

    toast.add({
      summary: 'Categoría creada',
      detail: 'La categoría se ha creado correctamente.',
      life: 3000,
      severity: 'success',
    });
  }
</script>

<template>
  <div class="flex justify-center">
    <Button size="small" @click="open = !open">
      <Plus />
      Nueva categoría
    </Button>

    <Dialog
      v-model:visible="open"
      header="Crear categoría"
      modal
      dismissableMask
      :style="{ width: '25rem' }"
    >
      <div class="flex flex-col gap-6">
        <span class="text-surface-500 dark:text-surface-400">
          Crea una nueva categoría para organizar tus productos.
        </span>

        <div class="flex flex-col gap-2">
          <InputText placeholder="Nombre de la categoría" v-model="form.name" />
          <Message v-if="error" severity="error" size="small" variant="simple">
            {{ error }}
          </Message>
        </div>
        <div class="flex justify-end gap-2">
          <Button label="Cancelar" severity="secondary" @click="open = false" />
          <Button label="Crear" @click="handleCreate" />
        </div>
      </div>
    </Dialog>
  </div>
</template>
