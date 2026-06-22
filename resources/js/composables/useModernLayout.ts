import { ref, watch, onMounted } from 'vue';

const MODERN_LAYOUT_KEY = 'is-modern-layout';

const isModern = ref<boolean>(
    typeof window !== 'undefined' ? localStorage.getItem(MODERN_LAYOUT_KEY) === 'true' : false
);

export function useModernLayout() {
    const toggleModern = (value: boolean) => {
        isModern.value = value;
        localStorage.setItem(MODERN_LAYOUT_KEY, String(value));
        updateBodyClass(value);
    };

    const updateBodyClass = (value: boolean) => {
        if (typeof document !== 'undefined') {
            document.documentElement.classList.toggle('modern-layout', value);
        }
    };

    watch(isModern, (value) => {
        updateBodyClass(value);
    }, { immediate: true });

    onMounted(() => {
        updateBodyClass(isModern.value);
    });

    return {
        isModern,
        toggleModern,
    };
}
