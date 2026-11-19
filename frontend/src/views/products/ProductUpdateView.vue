<script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { toast } from 'vue-sonner';
  import { Box } from '@lucide/vue';
  import ErrorAlert from '@/components/common/ErrorAlert.vue';
  import PageHeader from '@/components/PageHeader.vue';
  import { productService } from '@/services/productService';
  import type { ProductFormValues } from '@/types/product';
  import type { CategoryOptionResponse } from 'bindings/inventario/internal/modules/category/adapters/handler/models';
  import ProductForm from './ProductForm.vue';

  const route = useRoute();
  const router = useRouter();

  const categoryOptions = ref<CategoryOptionResponse[]>([]);
  const error = ref<string | null>(null);
  const loading = ref(false);
  const notFound = ref(false);

  const form = ref<ProductFormValues>({
    name: '',
    code: '',
    barcode: '',
    category_id: null,
    price: 0,
    min_stock: 0,
  });

  onMounted(async () => {
    const id = Number(route.params.id);

    const [product, opts] = await Promise.allSettled([
      productService.get(id),
      productService.categories(),
    ]);

    if (opts.status === 'fulfilled') {
      categoryOptions.value = opts.value ?? [];
    }

    if (product.status === 'rejected' || !product.value) {
      notFound.value = true;
      return;
    }

    const p = product.value;
    form.value = {
      name: p.name,
      code: p.code,
      barcode: p.barcode,
      category_id: p.category_id,
      price: p.price,
      min_stock: p.min_stock,
    };
  });

  async function handleUpdate() {
    error.value = null;
    loading.value = true;
    const id = Number(route.params.id);
    try {
      const { name, code, barcode, category_id, price, min_stock } = form.value;
      await productService.update(id, { name, code, barcode, category_id, price, min_stock });
      toast.success('Producto actualizado');
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
    <PageHeader title="Editar producto" :icon="Box" />

    <ErrorAlert v-if="notFound" message="Producto no encontrado." />

    <template v-else>
      <ErrorAlert v-if="error" :message="error" closable @close="error = null" />

      <ProductForm
        v-model="form"
        :category-options="categoryOptions"
        submit-label="Guardar"
        :loading="loading"
        @submit="handleUpdate"
        @cancel="router.back()"
      />
    </template>
  </div>
</template>
