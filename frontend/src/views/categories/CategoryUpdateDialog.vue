<script setup lang="ts">
import PencilLine from '@primeicons/vue/pencil';
import { ref, watch } from 'vue';
import { useCategories } from '@/composables/useCategories';
import type { CategoryResponse } from 'bindings/inventario/internal/handler/models';
import { useToast } from 'primevue/usetoast';

interface Props {
  category: CategoryResponse;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'updated'): void;
}>();

const { update, error, clearError } = useCategories();

const open = ref(false);
const name = ref('');

const toast = useToast();

watch(open, (val) => {
  if (val) {
    name.value = props.category.name;
    clearError();
  }
});

async function handleUpdate() {
  const success = await update(props.category.id, name.value);
  if (!success) return;
  open.value = false;
  emit('updated');

  toast.add({
    summary: 'Categoría actualizada',
    detail: 'La categoría se ha actualizado correctamente.',
    life: 3000,
    severity: 'success',
  });
}
</script>

<template>
  <Button
    severity="secondary"
    text
    rounded
    size="small"
    v-tooltip.top="'Editar'"
    @click="open = true"
  >
    <PencilLine :size="15" />
  </Button>

  <Dialog
    v-model:visible="open"
    header="Editar categoría"
    modal
    dismissableMask
    :style="{ width: '25rem' }"
  >
    <div class="flex flex-col gap-6">
      <span class="text-surface-500 dark:text-surface-400">
        Modifica el nombre de la categoría.
      </span>

      <div class="flex flex-col gap-2">
        <InputText placeholder="Nombre de la categoría" v-model="name" />
        <Message v-if="error" severity="error" size="small" variant="simple">
          {{ error }}
        </Message>
      </div>

      <div class="flex justify-end gap-2">
        <Button label="Cancelar" severity="secondary" @click="open = false" />
        <Button label="Guardar" @click="handleUpdate" />
      </div>
    </div>
  </Dialog>
</template>
