import React, { useState } from 'react';
import { X, Save } from 'lucide-react';

const VocabForm = ({ vocab, onSave, onCancel, selectedDate }) => {
  const [formData, setFormData] = useState({
    word: vocab?.word || '',
    meaning: vocab?.meaning || '',
    example: vocab?.example || '',
    status: vocab?.status || 'review_needed',
  });

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  const validateForm = () => {
    const newErrors = {};
    
    if (!formData.word.trim()) {
      newErrors.word = 'Word is required';
    }
    
    if (!formData.meaning.trim()) {
      newErrors.meaning = 'Meaning is required';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }
    
    setLoading(true);
    
    try {
      const dataToSave = {
        ...formData,
        date: selectedDate,
      };
      
      await onSave(dataToSave);
    } catch (error) {
      console.error('Error saving vocabulary:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="card max-w-md w-full p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-900">
            {vocab ? 'Edit Vocabulary' : 'Add New Vocabulary'}
          </h3>
          <button
            onClick={onCancel}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="word" className="block text-sm font-medium text-gray-700 mb-1">
              Word *
            </label>
            <input
              type="text"
              id="word"
              name="word"
              value={formData.word}
              onChange={handleChange}
              className={`input ${errors.word ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
              placeholder="Enter the word"
            />
            {errors.word && (
              <p className="mt-1 text-sm text-red-600">{errors.word}</p>
            )}
          </div>

          <div>
            <label htmlFor="meaning" className="block text-sm font-medium text-gray-700 mb-1">
              Meaning *
            </label>
            <textarea
              id="meaning"
              name="meaning"
              rows={3}
              value={formData.meaning}
              onChange={handleChange}
              className={`input ${errors.meaning ? 'border-red-300 focus:ring-red-500 focus:border-red-500' : ''}`}
              placeholder="Enter the meaning"
            />
            {errors.meaning && (
              <p className="mt-1 text-sm text-red-600">{errors.meaning}</p>
            )}
          </div>

          <div>
            <label htmlFor="example" className="block text-sm font-medium text-gray-700 mb-1">
              Example Sentence
            </label>
            <textarea
              id="example"
              name="example"
              rows={2}
              value={formData.example}
              onChange={handleChange}
              className="input"
              placeholder="Enter an example sentence (optional)"
            />
          </div>

          <div>
            <label htmlFor="status" className="block text-sm font-medium text-gray-700 mb-1">
              Status
            </label>
            <select
              id="status"
              name="status"
              value={formData.status}
              onChange={handleChange}
              className="input"
            >
              <option value="review_needed">Review Needed</option>
              <option value="mastered">Mastered</option>
            </select>
          </div>

          <div className="flex space-x-3 pt-4">
            <button
              type="submit"
              disabled={loading}
              className="btn btn-primary flex items-center space-x-2 flex-1"
            >
              <Save className="w-4 h-4" />
              <span>{loading ? 'Saving...' : 'Save'}</span>
            </button>
            <button
              type="button"
              onClick={onCancel}
              className="btn btn-secondary flex-1"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default VocabForm;