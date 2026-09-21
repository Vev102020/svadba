import { ref } from 'vue'

export function useQuest() {
    const questFiles = ref([])
    const questModalOpen = ref(false)
    const currentQuestId = ref(null)
    const isUploading = ref(false)
    const uploadProgress = ref(0)
    const showSuccess = ref(false)

    async function loadQuestFiles() {
        try {
            const res = await fetch('/api/quest/files', { cache: 'no-store' })
            if (!res.ok) throw new Error(`HTTP ${res.status}`)

            let data = await res.json()
            questFiles.value = (data || []).map(f => ({
                ...f,
                url: f.url + '?t=' + new Date().getTime()   // ← строчные url
            }))
        } catch (err) {
            console.error('[useQuest] Ошибка загрузки списка:', err)
            questFiles.value = [] 
        }
    }

    function openQuestModal(id) {
        currentQuestId.value = id
        questModalOpen.value = true
        showSuccess.value = false
        isUploading.value = false
        uploadProgress.value = 0
    }

    function closeQuestModal() {
        questModalOpen.value = false
        currentQuestId.value = null
    }

    function uploadQuest(file) {
        if (!currentQuestId.value) {
            console.error('[useQuest] Нет currentQuestId!')
            return
        }

        if (!file || !(file instanceof File)) {
            console.error('[useQuest] Некорректный файл:', file)
            alert('Выберите файл для загрузки')
            return
        }

        console.log('[useQuest] Загрузка файла для квеста', currentQuestId.value, file.name)

        const formData = new FormData()
        formData.append('file', file)

        isUploading.value = true
        uploadProgress.value = 0

        const xhr = new XMLHttpRequest()

        xhr.upload.addEventListener('progress', (e) => {
            if (e.lengthComputable) {
                uploadProgress.value = Math.round((e.loaded / e.total) * 100)
            }
        })

        xhr.addEventListener('load', () => {
            isUploading.value = false
            console.log('[useQuest] Ответ сервера:', xhr.status, xhr.responseText)

            if (xhr.status === 200) {
                showSuccess.value = true
                loadQuestFiles()
                setTimeout(() => {
                    closeQuestModal()
                }, 1500)
            } else {
                alert('Ошибка загрузки. Смотри консоль (F12).')
            }
        })

        xhr.addEventListener('error', () => {
            isUploading.value = false
            alert('Ошибка сети')
        })

        const url = `/api/quest/upload?id=${currentQuestId.value}`
        console.log('[useQuest] URL:', url)
        xhr.open('POST', url)
        xhr.send(formData)
    }

    return {
        questFiles,
        questModalOpen,
        currentQuestId,
        isUploading,
        uploadProgress,
        showSuccess,
        loadQuestFiles,
        openQuestModal,
        closeQuestModal,
        uploadQuest
    }
}