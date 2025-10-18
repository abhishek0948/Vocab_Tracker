import React, { useState, useEffect } from 'react';
import { format } from 'date-fns';
import { Plus, LogOut, BookOpen, Target, TrendingUp } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';
import { vocabAPI } from '../services/api';
import Calendar from '../components/Calendar';
import VocabList from '../components/VocabList';
import VocabForm from '../components/VocabForm';

const Dashboard = () => {
  const { user, logout } = useAuth();
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [vocabularies, setVocabularies] = useState([]);
  const [vocabCounts, setVocabCounts] = useState({});
  const [loading, setLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [editingVocab, setEditingVocab] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [stats, setStats] = useState({
    total: 0,
    mastered: 0,
    reviewNeeded: 0,
  });

  useEffect(() => {
    if (selectedDate) {
      fetchVocabulary();
    }
  }, [selectedDate, searchTerm]);

  useEffect(() => {
    fetchVocabCounts();
  }, []);

  const fetchVocabulary = async () => {
    setLoading(true);
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const response = await vocabAPI.getVocabulary(dateStr, searchTerm);
      setVocabularies(response.data.vocabularies);
      
      // Calculate stats
      const vocab = response.data.vocabularies;
      const mastered = vocab.filter(v => v.status === 'mastered').length;
      const reviewNeeded = vocab.filter(v => v.status === 'review_needed').length;
      
      setStats({
        total: vocab.length,
        mastered,
        reviewNeeded,
      });
    } catch (error) {
      console.error('Error fetching vocabulary:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchVocabCounts = async () => {
    try {
      // Fetch vocabulary for current month to get counts
      const response = await vocabAPI.getVocabulary('', '');
      const allVocab = response.data.vocabularies;
      
      // Group by date
      const counts = {};
      allVocab.forEach(vocab => {
        const date = vocab.date.split('T')[0]; // Extract date part
        counts[date] = (counts[date] || 0) + 1;
      });
      
      setVocabCounts(counts);
    } catch (error) {
      console.error('Error fetching vocabulary counts:', error);
    }
  };

  const handleSaveVocab = async (vocabData) => {
    try {
      if (editingVocab) {
        await vocabAPI.updateVocabulary(editingVocab.id, vocabData);
      } else {
        const dateStr = format(selectedDate, 'yyyy-MM-dd');
        await vocabAPI.createVocabulary({ ...vocabData, date: dateStr });
      }
      
      setShowForm(false);
      setEditingVocab(null);
      fetchVocabulary();
      fetchVocabCounts();
    } catch (error) {
      console.error('Error saving vocabulary:', error);
    }
  };

  const handleEditVocab = (vocab) => {
    setEditingVocab(vocab);
    setShowForm(true);
  };

  const handleDeleteVocab = async (id) => {
    try {
      await vocabAPI.deleteVocabulary(id);
      fetchVocabulary();
      fetchVocabCounts();
    } catch (error) {
      console.error('Error deleting vocabulary:', error);
    }
  };

  const handleUpdateStatus = async (id, newStatus) => {
    try {
      await vocabAPI.updateVocabulary(id, { status: newStatus });
      fetchVocabulary();
    } catch (error) {
      console.error('Error updating status:', error);
    }
  };

  const handleDateSelect = (date) => {
    setSelectedDate(date);
    setSearchTerm(''); // Clear search when changing dates
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-3">
              <BookOpen className="w-8 h-8 text-primary-600" />
              <h1 className="text-xl font-bold text-gray-900">Vocabulary Tracker</h1>
            </div>
            
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600">Welcome, {user?.email}</span>
              <button
                onClick={logout}
                className="btn btn-secondary flex items-center space-x-2"
              >
                <LogOut className="w-4 h-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Calendar */}
          <div className="lg:col-span-1">
            <Calendar
              selectedDate={selectedDate}
              onDateSelect={handleDateSelect}
              vocabCounts={vocabCounts}
            />
            
            {/* Stats */}
            <div className="mt-6 grid grid-cols-1 gap-4">
              <div className="card p-4">
                <div className="flex items-center space-x-3">
                  <div className="bg-blue-100 p-2 rounded-md">
                    <BookOpen className="w-5 h-5 text-blue-600" />
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Total Words</p>
                    <p className="text-xl font-semibold text-gray-900">{stats.total}</p>
                  </div>
                </div>
              </div>
              
              <div className="card p-4">
                <div className="flex items-center space-x-3">
                  <div className="bg-green-100 p-2 rounded-md">
                    <Target className="w-5 h-5 text-green-600" />
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Mastered</p>
                    <p className="text-xl font-semibold text-gray-900">{stats.mastered}</p>
                  </div>
                </div>
              </div>
              
              <div className="card p-4">
                <div className="flex items-center space-x-3">
                  <div className="bg-orange-100 p-2 rounded-md">
                    <TrendingUp className="w-5 h-5 text-orange-600" />
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Review Needed</p>
                    <p className="text-xl font-semibold text-gray-900">{stats.reviewNeeded}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Vocabulary List */}
          <div className="lg:col-span-2">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold text-gray-900">
                {format(selectedDate, 'MMMM d, yyyy')}
              </h2>
              <button
                onClick={() => setShowForm(true)}
                className="btn btn-primary flex items-center space-x-2"
              >
                <Plus className="w-4 h-4" />
                <span>Add Vocabulary</span>
              </button>
            </div>

            <VocabList
              vocabularies={vocabularies}
              onEdit={handleEditVocab}
              onDelete={handleDeleteVocab}
              onUpdateStatus={handleUpdateStatus}
              loading={loading}
              searchTerm={searchTerm}
              onSearchChange={setSearchTerm}
            />
          </div>
        </div>
      </div>

      {/* Vocabulary Form Modal */}
      {showForm && (
        <VocabForm
          vocab={editingVocab}
          selectedDate={format(selectedDate, 'yyyy-MM-dd')}
          onSave={handleSaveVocab}
          onCancel={() => {
            setShowForm(false);
            setEditingVocab(null);
          }}
        />
      )}
    </div>
  );
};

export default Dashboard;