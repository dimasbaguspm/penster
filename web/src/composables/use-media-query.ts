import { ref, onMounted, onUnmounted } from "vue";

export function useMediaQuery(query: string) {
  const matches = ref(false);
  let mediaQuery: MediaQueryList | null = null;

  const updateMatch = () => {
    if (mediaQuery) {
      matches.value = mediaQuery.matches;
    }
  };

  onMounted(() => {
    mediaQuery = window.matchMedia(query);
    updateMatch();
    mediaQuery.addEventListener("change", updateMatch);
  });

  onUnmounted(() => {
    if (mediaQuery) {
      mediaQuery.removeEventListener("change", updateMatch);
    }
  });

  return matches;
}

export function useIsMobile() {
  return useMediaQuery("(max-width: 767px)");
}

export function useIsDesktop() {
  return useMediaQuery("(min-width: 768px)");
}