<script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { toast } from 'vue-sonner';
  import { Box } from '@lucide/vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import PageHeader from '@/components/PageHeader.vue';
  import { productService } from '@/services/productService';
  import type { ProductFormValues } from '@/types/product';
  import type { CategoryOptionResponse } from 'bindings/inventario/internal/modules/category/adapters/handler/models';
  import ProductForm from './ProductForm.vue';

  const router = useRouter();

  const categoryOptions = ref<CategoryOptionResponse[]>([]);
  const error = ref<string | null>(null);
  const loading = ref(false);

  const form = ref<Required<ProductFormValues>>({
    name: '',
    code: '',
    barcode: '',
    category_id: null,
    price: 0,
    stock: 0,
    min_stock: 0,
  });

  onMounted(async () => {
    try {
      const opts = await productService.categories();
      categoryOptions.value = opts ?? [];
    } catch {
      // categories are optional
    }
  });

  async function handleCreate() {
    error.value = null;
    loading.value = true;
    try {
      await productService.create({ ...form.value });
      toast.success('Producto creado');
      router.push('/products');
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : String(cause);
    } finally {
      loading.value = false;
    }
  }
</script>

<template>
  <div class="flex max-w-xl flex-col gap-6">
    <PageHeader title="Nuevo producto" :icon="Box" />

    <ErrorAlert v-if="error" :message="error" closable @close="error = null" />

    <ProductForm
      v-model="form"
      :category-options="categoryOptions"
      submit-label="Crear"
      :loading="loading"
      @submit="handleCreate"
      @cancel="router.back()"
    />
  </div>
</template>
