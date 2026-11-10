<script setup lang="ts">
  import Refresh from '@primeicons/vue/refresh';
  import { useConfirm } from 'primevue/useconfirm';
  import { useToast } from 'primevue/usetoast';
  import { useCategories } from '@/composables/useCategories';
  import type { Category } from '@/types/category';
  import { InfoCircle } from '@primeicons/vue';

  interface Props {
    category: Category;
  }

  const props = defineProps<Props>();

  const emit = defineEmits<{
    (e: 'restored'): void;
  }>();

  const { restore } = useCategories();
  const confirm = useConfirm();
  const toast = useToast();

  function handleRestore() {
    confirm.require({
      message: `¿Restaurar la categoría "${props.category.name}"?`,
      header: 'Confirmar restauración',
      icon: InfoCircle,
      rejectProps: {
        label: 'Cancel',
        severity: 'secondary',
        outlined: true,
      },
      acceptProps: {
        label: 'Restaurar',
        severity: 'success',
      },
      accept: async () => {
        const ok = await restore(props.category.id);
        if (ok) {
          toast.add({ severity: 'success', summary: 'Categoría restaurada', life: 3000 });
          emit('restored');
        }
      },
    });
  }
</script>

<template>
  <Button
    severity="success"
    text
    rounded
    size="small"
    v-tooltip.top="'Restaurar'"
    @click="handleRestore"
  >
    <Refresh :size="15" />
  </Button>
</template>
