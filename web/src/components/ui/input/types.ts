export type InputType = "text" | "email" | "password" | "number" | "search" | "date";

export interface InputProps {
  label?: string;
  type?: InputType;
  placeholder?: string;
  modelValue?: string | number;
  disabled?: boolean;
  error?: string;
}

export interface InputEmits {
  "update:modelValue": [value: string];
}
