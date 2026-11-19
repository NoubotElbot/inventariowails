<script setup lang="ts">
  import { computed } from 'vue';
  import AppSelect from '@/components/common/AppSelect.vue';
  import FormField from '@/components/common/FormField.vue';
  import NumberInput from '@/components/common/NumberInput.vue';
  import { toSelectOptions } from '@/components/common/select';
  import { Button } from '@/components/ui/button';
  import { Input } from '@/components/ui/input';
  import { Spinner } from '@/components/ui/spinner';
  import type { ProductFormValues } from '@/types/product';
  import type { CategoryOptionResponse } from 'bindings/inventario/internal/modules/category/adapters/handler/models';

  /** Shared create/edit product form. The stock field shows only when `form.stock` is set. */
  const form = defineModel<ProductFormValues>({ required: true });

  const {
    categoryOptions,
    submitLabel,
    loading = false,
  } = defineProps<{
    categoryOptions: readonly CategoryOptionResponse[];
    submitLabel: string;
    loading?: boolean;
  }>();

  const emit = defineEmits<{
    submit: [];
    cancel: [];
  }>();

  const options = computed(() => toSelectOptions(categoryOptions));
  const withStock = computed(() => form.value.stock !== undefined);
</script>

<template>
  <form class="flex flex-col gap-6" @submit.prevent="emit('submit')">
    <div class="flex flex-col gap-4">
      <FormField label="Nombre" for="product-name" required>
        <Input id="product-name" v-model="form.name" placeholder="Nombre del producto" />
      </FormField>

      <div class="grid grid-cols-2 gap-4">
        <FormField label="Código" for="product-code" required>
          <Input id="product-code" v-model="form.code" placeholder="Ej. PRD-001" />
        </FormField>
        <FormField label="Código de barras" for="product-barcode">
          <Input id="product-barcode" v-model="form.barcode" placeholder="Ej. 7501234567890" />
        </FormField>
      </div>

      <FormField label="Categoría" for="product-category">
        <AppSelect
          id="product-category"
          v-model="form.category_id"
          :options="options"
          placeholder="Sin categoría"
          clearable
        />
      </FormField>

      <div class="grid gap-4" :class="withStock ? 'grid-cols-3' : 'grid-cols-2'">
        <FormField label="Precio" for="product-price" required>
          <NumberInput id="product-price" v-model="form.price" :min="0" currency />
        </FormField>
        <FormField v-if="form.stock !== undefined" label="Stock inicial" for="product-stock">
          <NumberInput id="product-stock" v-model="form.stock" :min="0" />
        </FormField>
        <FormField label="Stock mínimo" for="product-min-stock">
          <NumberInput id="product-min-stock" v-model="form.min_stock" :min="0" />
        </FormField>
      </div>
    </div>

    <div class="flex justify-end gap-2">
      <Button type="button" variant="outline" @click="emit('cancel')">Cancelar</Button>
      <Button type="submit" :disabled="loading">
        <Spinner v-if="loading" />
        {{ submitLabel }}
      </Button>
    </div>
  </form>
</template>
