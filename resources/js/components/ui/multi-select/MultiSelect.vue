<script setup lang="ts">
import { ref, computed } from 'vue';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
    CommandSeparator,
} from '@/components/ui/command';
import {
    Drawer,
    DrawerContent,
    DrawerHeader,
    DrawerTitle,
    DrawerTrigger,
} from "@/components/ui/drawer";
import { Check, ChevronsUpDown, X } from 'lucide-vue-next';
import { cn } from '@/lib/utils';

interface Option {
    label: string;
    value: string | number;
}

interface GroupedOption {
    group: string;
    options: Option[];
}

type MultiSelectOption = Option | GroupedOption;

const props = withDefaults(defineProps<{
    options: MultiSelectOption[];
    modelValue: (string | number)[];
    placeholder?: string;
    searchPlaceholder?: string;
    emptyMessage?: string;
    disabled?: boolean;
    shouldScaleBackground?: boolean;
}>(), {
    shouldScaleBackground: true,
});

const emit = defineEmits(['update:modelValue']);

const open = ref(false);

const flatOptions = computed(() => {
    return props.options.reduce((acc: Option[], item) => {
        if ('group' in item) {
            return [...acc, ...item.options];
        }
        return [...acc, item];
    }, []);
});

function toggleOption(value: string | number) {
    const newValue = [...props.modelValue];
    const index = newValue.indexOf(value);
    if (index > -1) {
        newValue.splice(index, 1);
    } else {
        newValue.push(value);
    }
    emit('update:modelValue', newValue);
}

function removeOption(value: string | number) {
    if (props.disabled) return;
    const newValue = props.modelValue.filter((val) => val !== value);
    emit('update:modelValue', newValue);
}
</script>

<template>
    <Drawer :open="disabled ? false : open" @update:open="open = $event"
        :should-scale-background="shouldScaleBackground">
        <DrawerTrigger as-child>
            <Button variant="outline" role="combobox" :aria-expanded="open" :disabled="disabled"
                class="w-full justify-between min-h-11 h-auto px-3 py-2 hover:bg-background border-input font-normal disabled:opacity-50 disabled:cursor-not-allowed text-right"
                dir="rtl">
                <div class="flex flex-wrap gap-1 items-center overflow-hidden">
                    <template v-if="modelValue.length > 0">
                        <Badge v-for="val in modelValue" :key="val" variant="secondary" :class="cn(
                            'rounded-sm px-1 font-normal flex items-center gap-1 bg-primary/10 text-primary border-primary/20',
                            disabled && 'opacity-70 pointer-events-none'
                        )" @click.stop="removeOption(val)">
                            {{flatOptions.find(o => o.value === val)?.label}}
                            <X v-if="!disabled" class="h-3 w-3 hover:text-destructive cursor-pointer" />
                        </Badge>
                    </template>
                    <span v-else class="text-muted-foreground">{{ placeholder || 'اختر الخيارات...' }}</span>
                </div>
                <ChevronsUpDown class="mr-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
        </DrawerTrigger>
        <DrawerContent class="p-0 h-[90vh] max-h-[96vh] md:max-w-3xl md:mx-auto">
            <DrawerHeader class="border-b px-4 py-3">
                <DrawerTitle>{{ placeholder || 'اختر الخيارات' }}</DrawerTitle>
            </DrawerHeader>
            <div class="overflow-y-auto flex-1">
                <Command dir="rtl" class="h-full">
                    <CommandInput :placeholder="searchPlaceholder || 'ابحث...'" />
                    <CommandEmpty>{{ emptyMessage || 'لا يوجد نتائج.' }}</CommandEmpty>
                    <CommandList class="max-h-none h-full">
                        <template v-for="(item, idx) in options" :key="idx">
                            <CommandGroup v-if="'group' in item" :heading="item.group">
                                <CommandItem v-for="option in item.options" :key="option.value"
                                    :value="String(option.label)" @select="toggleOption(option.value)"
                                    class="cursor-pointer py-3">
                                    <div :class="cn(
                                        'ml-2 flex h-5 w-5 items-center justify-center rounded-sm border border-primary',
                                        modelValue.includes(option.value)
                                            ? 'bg-primary text-primary-foreground'
                                            : 'opacity-50 [&_svg]:invisible'
                                    )">
                                        <Check class="h-4 w-4" />
                                    </div>
                                    <span class="flex-1 text-sm">{{ option.label }}</span>
                                    <code class="text-[10px] opacity-40 ml-2 font-mono">{{ option.value }}</code>
                                </CommandItem>
                            </CommandGroup>

                            <CommandGroup v-else>
                                <CommandItem :key="item.value" :value="String(item.label)"
                                    @select="toggleOption(item.value)" class="cursor-pointer py-3">
                                    <div :class="cn(
                                        'ml-2 flex h-5 w-5 items-center justify-center rounded-sm border border-primary',
                                        modelValue.includes(item.value)
                                            ? 'bg-primary text-primary-foreground'
                                            : 'opacity-50 [&_svg]:invisible'
                                    )">
                                        <Check class="h-4 w-4" />
                                    </div>
                                    <span class="flex-1 text-sm">{{ item.label }}</span>
                                </CommandItem>
                            </CommandGroup>
                            <CommandSeparator v-if="idx < options.length - 1" />
                        </template>
                    </CommandList>
                </Command>
            </div>
        </DrawerContent>
    </Drawer>
</template>
