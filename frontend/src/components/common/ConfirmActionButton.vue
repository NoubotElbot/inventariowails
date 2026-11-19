<script setup lang="ts">
  import type { Component } from 'vue';
  import { shallowRef } from 'vue';
  import { toast } from 'vue-sonner';
  import IconButton from '@/components/common/IconButton.vue';
  import { useConfirm, type ConfirmOptions } from '@/composables/useConfirm';

  /**
   * Icon button that asks for confirmation, runs `action` and reports the
   * outcome with a toast. Used for delete / restore row actions.
   */
  const {
    icon,
    tooltip,
    confirm: confirmOptions,
    action,
    successMessage,
    variant = 'destructive',
  } = defineProps<{
    icon: Component;
    tooltip: string;
    confirm: ConfirmOptions;
    /** Resolves to `true` on success; the caller surfaces its own errors. */
    action: () => Promise<boolean>;
    successMessage: string;
    variant?: 'destructive' | 'success';
  }>();

  const emit = defineEmits<{
    done: [];
  }>();

  const { confirm } = useConfirm();
  const running = shallowRef(false);

  async function run(): Promise<void> {
    const accepted = await confirm({ variant, ...confirmOptions });
    if (!accepted) return;

    running.value = true;
    try {
      if (!(await action())) return;
      toast.success(successMessage);
      emit('done');
    } finally {
      running.value = false;
    }
  }
</script>

<template>
  <IconButton
    :icon="icon"
    :tooltip="tooltip"
    :class="
      variant === 'destructive'
        ? 'text-destructive hover:text-destructive'
        : 'text-success hover:text-success'
    "
    :disabled="running"
    @click="run"
  />
</template>
