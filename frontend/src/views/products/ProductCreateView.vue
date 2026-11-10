<script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { useToast } from 'primevue/usetoast';
  import Box from '@primeicons/vue/box';
  import { productService } from '@/services/productService';
  import type { CategoryOptionResponse } from 'bindings/inventario/internal/modules/category/handler/models';

  const router = useRouter();
  const toast = useToast();

  const categoryOptions = ref<CategoryOptionResponse[]>([]);
  const error = ref<string | null>(null);
  const loading = ref(false);

  const form = ref({
    name: '',
    code: '',
    barcode: '',
    category_id: null as number | null,
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
      await productService.create({
        name: form.value.name,
        code: form.value.code,
        barcode: form.value.barcode,
        category_id: form.value.category_id,
        price: form.value.price,
        stock: form.value.stock,
        min_stock: form.value.min_stock,
      });
      toast.add({ severity: 'success', summary: 'Producto creado', life: 3000 });
      router.push('/products');
    } catch (cause) {
      error.value = cause instanceof Error ? cause.message : String(cause);
    } finally {
      loading.value = false;
    }
  }
</script>

<template>
  <div class="flex flex-col gap-6 max-w-xl">
    <!-- Encabezado -->
    <div class="flex items-center gap-3">
      <Box :size="20" class="text-primary" />
      <h1 class="text-xl font-semibold text-surface-800 dark:text-surface-100">Nuevo producto</h1>
    </div>

    <Message v-if="error" severity="error" :closable="true" @close="error = null">
      {{ error }}
    </Message>

    <div class="flex flex-col gap-4">
      <div class="flex flex-col gap-1">
        <label class="text-sm font-medium text-surface-700 dark:text-surface-300">Nombre *</label>
        <InputText v-model="form.name" placeholder="Nombre del producto" />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-surface-700 dark:text-surface-300">Código *</label>
          <InputText v-model="form.code" placeholder="Ej. PRD-001" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >Código de barras</label
          >
          <InputText v-model="form.barcode" placeholder="Ej. 7501234567890" />
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm font-medium text-surface-700 dark:text-surface-300">Categoría</label>
        <Select
          v-model="form.category_id"
          :options="categoryOptions"
          option-label="name"
          option-value="id"
          placeholder="Sin categoría"
          show-clear
        />
      </div>

      <div class="grid grid-cols-3 gap-4">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-surface-700 dark:text-surface-300">Precio *</label>
          <InputNumber
            v-model="form.price"
            mode="currency"
            currency="CLP"
            locale="es-CL"
            :min="0"
          />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >Stock inicial</label
          >
          <InputNumber v-model="form.stock" :min="0" />
        </div>
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-surface-700 dark:text-surface-300"
            >Stock mínimo</label
          >
          <InputNumber v-model="form.min_stock" :min="0" />
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-2">
      <Button label="Cancelar" severity="secondary" @click="router.back()" />
      <Button label="Crear" :loading="loading" @click="handleCreate" />
    </div>
  </div>
</template>
