<script setup lang="ts">
  import { Check } from '@lucide/vue';
  import { Button } from '@/components/ui/button';
  import { Spinner } from '@/components/ui/spinner';
  import { formatPrice } from '@/utils/currency';

  defineProps<{
    units: number;
    total: number;
    loading: boolean;
    disabled: boolean;
  }>();

  const emit = defineEmits<{
    cancel: [];
    confirm: [];
  }>();
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-4 rounded-lg border px-4 py-3">
    <div class="flex flex-col">
      <span class="text-sm text-muted-foreground">{{ units }} unidades</span>
      <span class="text-2xl font-semibold text-foreground">{{ formatPrice(total) }}</span>
    </div>

    <div class="flex gap-2">
      <Button variant="outline" :disabled="loading" @click="emit('cancel')">Cancelar</Button>
      <Button :disabled="disabled || loading" @click="emit('confirm')">
        <Spinner v-if="loading" />
        <Check v-else />
        Confirmar venta
      </Button>
    </div>
  </div>
</template>
