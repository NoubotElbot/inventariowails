<script setup lang="ts">
  import { computed, ref, watch } from 'vue';
  import { toast } from 'vue-sonner';
  import FormField from '@/components/common/FormField.vue';
  import { Button } from '@/components/ui/button';
  import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
  } from '@/components/ui/dialog';
  import { Input } from '@/components/ui/input';
  import { Spinner } from '@/components/ui/spinner';
  import { useCategories } from '@/composables/useCategories';
  import type { Category } from '@/types/category';

  /** Creates a category, or edits `category` when one is given. */
  const open = defineModel<boolean>('open', { required: true });

  const { category = null } = defineProps<{
    category?: Category | null;
  }>();

  const emit = defineEmits<{
    saved: [];
  }>();

  const { create, update, loading, error, clearError } = useCategories();

  const name = ref('');
  const isEdit = computed(() => category !== null);

  const copy = computed(() =>
    isEdit.value
      ? {
          title: 'Editar categoría',
          description: 'Modifica el nombre de la categoría.',
          submit: 'Guardar',
          success: 'Categoría actualizada',
        }
      : {
          title: 'Crear categoría',
          description: 'Crea una nueva categoría para organizar tus productos.',
          submit: 'Crear',
          success: 'Categoría creada',
        }
  );

  // Every time the dialog opens, start from the category (or blank) with no stale error.
  watch(open, (isOpen) => {
    if (!isOpen) return;
    name.value = category?.name ?? '';
    clearError();
  });

  async function submit(): Promise<void> {
    const ok = category ? await update(category.id, name.value) : await create(name.value);
    if (!ok) return;

    open.value = false;
    emit('saved');
    toast.success(copy.value.success);
  }
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-sm">
      <form class="contents" @submit.prevent="submit">
        <DialogHeader>
          <DialogTitle>{{ copy.title }}</DialogTitle>
          <DialogDescription>{{ copy.description }}</DialogDescription>
        </DialogHeader>

        <FormField label="Nombre" for="category-name" required :error="error">
          <Input
            id="category-name"
            v-model="name"
            placeholder="Nombre de la categoría"
            :aria-invalid="error ? true : undefined"
          />
        </FormField>

        <DialogFooter>
          <Button type="button" variant="outline" @click="open = false">Cancelar</Button>
          <Button type="submit" :disabled="loading">
            <Spinner v-if="loading" />
            {{ copy.submit }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
