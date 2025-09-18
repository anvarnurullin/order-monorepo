import React, { useState } from 'react';
import './ImageUpload.scss';

interface ImageUploadProps {
  productId: number;
  onImageUploaded: (imageUrl: string) => void;
}

export const ImageUpload: React.FC<ImageUploadProps> = ({ productId, onImageUploaded }) => {
    const [uploading, setUploading] = useState(false);

    async function uploadProductImage(productId: number, file: File): Promise<{ image_url: string }> {
        const formData = new FormData();
        formData.append('image', file);

        const res = await fetch(`/api/v1/products/${productId}/image`, {
            method: 'POST',
            body: formData,
        });

        if (!res.ok) {
            throw new Error('Failed to upload image');
        }

        return res.json();
    }

    const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        setUploading(true);
        try {
        const result = await uploadProductImage(productId, file);
        onImageUploaded(result.image_url);
        } catch (error) {
        console.error('Failed to upload image:', error);
        alert('Ошибка загрузки изображения');
        } finally {
        setUploading(false);
        }
    };

    return (
        <div className="image-upload">
        <input
            type="file"
            accept="image/*"
            onChange={handleFileChange}
            disabled={uploading}
            className="image-upload__input"
        />
        {uploading && <span className="image-upload__status">Загрузка...</span>}
        </div>
    );
};