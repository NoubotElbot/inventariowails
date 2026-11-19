<script setup lang="ts">
  import { computed, shallowRef, watch } from 'vue';
  import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
  } from '@/components/ui/alert-dialog';
  import { Button } from '@/components/ui/button';
  import { useConfirmHost } from '@/composables/useConfirm';

  /** Renders the requests made through `useConfirm()`. Mount it once, in App.vue. */
  const { request, settle } = useConfirmHost();

  // Keeps the last request on screen while the dialog plays its close animation.
  const shown = shallowRef(request.value);
  watch(request, (value) => {
    if (value) shown.value = value;
  });

  const open = computed({
    get: () => request.value !== null,
    // Closing by Escape counts as a cancellation.
    set: (value) => {
      if (!value) settle(false);
    },
  });
</script>

<template>
  <AlertDialog v-model:open="open">
    <AlertDialogContent v-if="shown">
      <AlertDialogHeader>
        <AlertDialogTitle>{{ shown.title }}</AlertDialogTitle>
        <AlertDialogDescription>{{ shown.description }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <Button variant="outline" @click="settle(false)">{{ shown.cancelLabel }}</Button>
        <Button :variant="shown.variant" @click="settle(true)">
          {{ shown.confirmLabel }}
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
