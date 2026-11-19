<script setup lang="ts">
  import { Field, FieldError, FieldLabel } from '@/components/ui/field';

  /** Label + control + optional error message. The control goes in the default slot. */
  const {
    label,
    for: htmlFor,
    required = false,
    error,
  } = defineProps<{
    label: string;
    /** Id of the control, so clicking the label focuses it. */
    for?: string;
    required?: boolean;
    error?: string | null;
  }>();
</script>

<template>
  <Field :data-invalid="error ? true : undefined">
    <FieldLabel :for="htmlFor">
      {{ label }}<span v-if="required" class="text-destructive" aria-hidden="true">*</span>
    </FieldLabel>
    <slot />
    <FieldError v-if="error">{{ error }}</FieldError>
  </Field>
</template>
