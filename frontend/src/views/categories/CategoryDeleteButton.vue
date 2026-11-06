<script setup lang="ts">
import Trash from '@primeicons/vue/trash';
import { useConfirm } from 'primevue/useconfirm';
import { useToast } from 'primevue/usetoast';
import { useCategories } from '@/composables/useCategories';
import type { CategoryResponse } from 'bindings/inventario/internal/handler/models';
import { InfoCircle } from '@primeicons/vue';

interface Props {
  category: CategoryResponse;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'deleted'): void;
}>();

const { remove } = useCategories();
const confirm = useConfirm();
const toast = useToast();

function handleDelete() {
  confirm.require({
    message: `¿Eliminar la categoría "${props.category.name}"?`,
    header: 'Confirmar eliminación',
    icon: InfoCircle,
    rejectProps: {
      label: 'Cancel',
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: 'Eliminar',
      severity: 'danger',
    },
    accept: async () => {
      const ok = await remove(props.category.id);
      if (ok) {
        toast.add({ severity: 'success', summary: 'Categoría eliminada', life: 3000 });
        emit('deleted');
      }
    },
  });
}
</script>

<template>
  <Button
    severity="danger"
    text
    rounded
    size="small"
    v-tooltip.top="'Eliminar'"
    @click="handleDelete"
  >
    <Trash :size="15" />
  </Button>
</template>
