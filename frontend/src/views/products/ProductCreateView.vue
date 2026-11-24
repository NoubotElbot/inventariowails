<script setup lang="ts">
  import { nextTick, onMounted, ref } from 'vue';
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

  function emptyForm(): Required<ProductFormValues> {
    return {
      name: '',
      code: '',
      barcode: '',
      category_id: null,
      price: 0,
      stock: 0,
      min_stock: 0,
    };
  }

  const form = ref(emptyForm());

  onMounted(async () => {
    try {
      const opts = await productService.categories();
      categoryOptions.value = opts ?? [];
    } catch {
      // categories are optional
    }
  });

  /** `keepGoing` stays on the page with a fresh form (category kept) to register another one. */
  async function handleCreate(keepGoing = false) {
    error.value = null;
    loading.value = true;
    try {
      await productService.create({ ...form.value });
      toast.success('Producto creado');
      if (keepGoing) {
        form.value = { ...emptyForm(), category_id: form.value.category_id };
        await nextTick();
        document.getElementById('product-name')?.focus();
      } else {
        router.push('/products');
      }
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
      continue-label="Crear y agregar otro"
      :loading="loading"
      @submit="handleCreate()"
      @submit-continue="handleCreate(true)"
      @cancel="router.back()"
    />
  </div>
</template>
