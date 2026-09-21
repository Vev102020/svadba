import { ref, computed } from 'vue'

export function useLightbox(files) {
    const lightboxOpen = ref(false)
    const lightboxIndex = ref(0)
    const lightboxFiles = ref([])

    const currentFile = computed(() => files.value[lightboxIndex.value] ?? null)

    function openLightbox(files, index) {
        lightboxFiles.value = files
        lightboxIndex.value = index
        lightboxOpen.value = true
        document.body.style.overflow = 'hidden'
    }

    function closeLightbox() {
        lightboxOpen.value = false
        document.body.style.overflow = ''
        document.querySelectorAll('.lightbox__content video').forEach(v => {
            v.pause()
            v.currentTime = 0
        })
    }

    function prev() {
        if (lightboxIndex.value > 0) lightboxIndex.value--
    }

    function next() {
        if (lightboxIndex.value < files.value.length - 1) lightboxIndex.value++
    }

    return {
        lightboxOpen,
        lightboxIndex,
        lightboxFiles,
        currentFile,
        openLightbox,
        closeLightbox,
        prev,
        next
    }
}