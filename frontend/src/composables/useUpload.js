import { ref } from 'vue'

export function useUpload(files, loadFiles) {
    const isUploading = ref(false)
    const uploadProgress = ref(0)
    const uploadModalOpen = ref(false)
    const selectedFiles = ref([])
    const showSuccess = ref(false)

    function openUploadModal() {
        uploadModalOpen.value = true
        showSuccess.value = false
    }

    function closeUploadModal() {
        uploadModalOpen.value = false
        selectedFiles.value.forEach(f => {
            if (f.preview) URL.revokeObjectURL(f.preview)
        })
        selectedFiles.value = []
        showSuccess.value = false
    }

    function addFiles(fileList) {
        Array.from(fileList).forEach(file => {
            const isImage = file.type.startsWith('image/')
            const preview = isImage ? URL.createObjectURL(file) : null
            selectedFiles.value.push({
                file,
                name: file.name,
                type: file.type,
                size: file.size,
                preview,
                isImage
            })
        })
    }

    function removeFile(index) {
        const item = selectedFiles.value[index]
        if (item.preview) URL.revokeObjectURL(item.preview)
        selectedFiles.value.splice(index, 1)
    }

    function uploadSelected() {
        if (selectedFiles.value.length === 0) return

        const formData = new FormData()
        selectedFiles.value.forEach(item => formData.append('files', item.file))

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
            if (xhr.status === 200) {
                showSuccess.value = true
                loadFiles()
                setTimeout(() => {
                    closeUploadModal()
                }, 2000)
            } else {
                alert('Ошибка при загрузке файлов')
            }
        })

        xhr.addEventListener('error', () => {
            isUploading.value = false
            alert('Ошибка сети. Проверьте подключение.')
        })

        xhr.open('POST', '/api/upload')
        xhr.send(formData)
    }

    return {
        isUploading,
        uploadProgress,
        uploadModalOpen,
        selectedFiles,
        showSuccess,
        openUploadModal,
        closeUploadModal,
        addFiles,
        removeFile,
        uploadSelected
    }
}