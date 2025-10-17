import React, { useState } from 'react';
import { Edit, Trash2, CheckCircle, Circle, Search } from 'lucide-react';

const VocabList = ({ vocabularies, onEdit, onDelete, onUpdateStatus, loading, searchTerm, onSearchChange }) => {
  const [deletingId, setDeletingId] = useState(null);

  const handleDelete = async (id) => {
    setDeletingId(id);
    await onDelete(id);
    setDeletingId(null);
  };

  const getStatusColor = (status) => {
    return status === 'mastered' 
      ? 'text-green-600 bg-green-100' 
      : 'text-orange-600 bg-orange-100';
  };

  const getStatusIcon = (status) => {
    return status === 'mastered' 
      ? <CheckCircle className="w-4 h-4" />
      : <Circle className="w-4 h-4" />;
  };

  if (loading) {
    return (
      <div className="card p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-4 bg-gray-200 rounded w-1/4"></div>
          <div className="space-y-3">
            {[1, 2, 3].map(i => (
              <div key={i} className="h-20 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="card p-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4 sm:mb-0">
          Vocabulary ({vocabularies.length})
        </h3>
        
        <div className="relative">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center">
            <Search className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="text"
            placeholder="Search words..."
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            className="input pl-10 w-full sm:w-64"
          />
        </div>
      </div>

      {vocabularies.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-2">
            <Circle className="w-12 h-12 mx-auto" />
          </div>
          <p className="text-gray-500">
            {searchTerm ? 'No vocabulary matches your search.' : 'No vocabulary entries for this date.'}
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {vocabularies.map((vocab) => (
            <div
              key={vocab.id}
              className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center space-x-3 mb-2">
                    <h4 className="text-lg font-semibold text-gray-900">
                      {vocab.word}
                    </h4>
                    <button
                      onClick={() => onUpdateStatus(vocab.id, vocab.status === 'mastered' ? 'review_needed' : 'mastered')}
                      className={`flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium transition-colors ${getStatusColor(vocab.status)}`}
                    >
                      {getStatusIcon(vocab.status)}
                      <span>{vocab.status === 'mastered' ? 'Mastered' : 'Review Needed'}</span>
                    </button>
                  </div>
                  
                  <p className="text-gray-700 mb-2">
                    <span className="font-medium">Meaning:</span> {vocab.meaning}
                  </p>
                  
                  {vocab.example && (
                    <p className="text-gray-600 text-sm italic">
                      <span className="font-medium">Example:</span> {vocab.example}
                    </p>
                  )}
                </div>

                <div className="flex space-x-2 ml-4">
                  <button
                    onClick={() => onEdit(vocab)}
                    className="p-2 text-gray-400 hover:text-primary-600 hover:bg-primary-50 rounded-md transition-colors"
                    title="Edit vocabulary"
                  >
                    <Edit className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => handleDelete(vocab.id)}
                    disabled={deletingId === vocab.id}
                    className="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors disabled:opacity-50"
                    title="Delete vocabulary"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default VocabList;