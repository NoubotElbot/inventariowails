<script setup lang="ts">
import Trash from '@primeicons/vue/trash';
import { InfoCircle } from '@primeicons/vue';
import { useConfirm } from 'primevue/useconfirm';
import { useToast } from 'primevue/usetoast';
import { useProducts } from '@/composables/useProducts';
import type { ProductResponse } from 'bindings/inventario/internal/handler/models';

interface Props {
  product: ProductResponse;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'deleted'): void;
}>();

const { remove } = useProducts();
const confirm = useConfirm();
const toast = useToast();

function handleDelete() {
  confirm.require({
    message: `¿Eliminar el producto "${props.product.name}"?`,
    header: 'Confirmar eliminación',
    icon: InfoCircle,
    rejectProps: {
      label: 'Cancelar',
      severity: 'secondary',
      outlined: true,
    },
    acceptProps: {
      label: 'Eliminar',
      severity: 'danger',
    },
    accept: async () => {
      const ok = await remove(props.product.id);
      if (ok) {
        toast.add({ severity: 'success', summary: 'Producto eliminado', life: 3000 });
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
