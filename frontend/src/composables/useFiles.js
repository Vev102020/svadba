import { ref } from 'vue'

export function useFiles() {
    const files = ref([])

    async function loadFiles() {
        try {
            const res = await fetch('/api/files', { cache: 'no-store' })
            if (!res.ok) throw new Error(`HTTP ${res.status}`)

            let data = await res.json()
            files.value = data || []
        } catch (err) {
            console.error('[useFiles] Ошибка загрузки списка:', err)
            files.value = []
        }
    }

    return { files, loadFiles }
}