<script setup lang="ts">
  import { CircleAlert, X } from '@lucide/vue';
  import { Alert, AlertAction, AlertDescription } from '@/components/ui/alert';
  import { Button } from '@/components/ui/button';

  /** Error banner. With `closable` it shows a dismiss button that emits `close`. */
  const { message, closable = false } = defineProps<{
    message: string;
    closable?: boolean;
  }>();

  defineEmits<{
    close: [];
  }>();
</script>

<template>
  <Alert variant="destructive">
    <CircleAlert />
    <AlertDescription>{{ message }}</AlertDescription>
    <AlertAction v-if="closable">
      <Button variant="ghost" size="icon-xs" aria-label="Cerrar" @click="$emit('close')">
        <X />
      </Button>
    </AlertAction>
  </Alert>
</template>
