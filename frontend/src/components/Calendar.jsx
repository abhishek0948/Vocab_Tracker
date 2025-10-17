import React, { useState, useEffect } from 'react';
import { format, startOfMonth, endOfMonth, eachDayOfInterval, getDay, isSameMonth, isSameDay, isToday } from 'date-fns';
import { ChevronLeft, ChevronRight } from 'lucide-react';

const Calendar = ({ selectedDate, onDateSelect, vocabCounts = {} }) => {
  const [currentMonth, setCurrentMonth] = useState(new Date());

  const monthStart = startOfMonth(currentMonth);
  const monthEnd = endOfMonth(currentMonth);
  const days = eachDayOfInterval({ start: monthStart, end: monthEnd });

  const startDay = getDay(monthStart);
  const emptyCells = Array(startDay).fill(null);

  const prevMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() - 1));
  };

  const nextMonth = () => {
    setCurrentMonth(new Date(currentMonth.getFullYear(), currentMonth.getMonth() + 1));
  };

  const getDayKey = (date) => format(date, 'yyyy-MM-dd');

  const getDayClasses = (day) => {
    let classes = 'relative w-10 h-10 flex items-center justify-center text-sm cursor-pointer rounded-md transition-colors ';
    
    if (!isSameMonth(day, currentMonth)) {
      classes += 'text-gray-300 ';
    } else if (isToday(day)) {
      classes += 'bg-primary-100 text-primary-800 font-semibold ';
    } else {
      classes += 'text-gray-700 hover:bg-gray-100 ';
    }
    
    if (selectedDate && isSameDay(day, selectedDate)) {
      classes += 'bg-primary-600 text-white hover:bg-primary-700 ';
    }
    
    return classes;
  };

  return (
    <div className="card p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900">
          {format(currentMonth, 'MMMM yyyy')}
        </h2>
        <div className="flex space-x-1">
          <button
            onClick={prevMonth}
            className="p-2 hover:bg-gray-100 rounded-md transition-colors"
          >
            <ChevronLeft className="w-5 h-5" />
          </button>
          <button
            onClick={nextMonth}
            className="p-2 hover:bg-gray-100 rounded-md transition-colors"
          >
            <ChevronRight className="w-5 h-5" />
          </button>
        </div>
      </div>

      <div className="grid grid-cols-7 gap-1 mb-2">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="text-center text-xs font-medium text-gray-500 py-2">
            {day}
          </div>
        ))}
      </div>

      <div className="grid grid-cols-7 gap-1">
        {emptyCells.map((_, index) => (
          <div key={`empty-${index}`} className="w-10 h-10"></div>
        ))}
        
        {days.map((day) => {
          const dayKey = getDayKey(day);
          const vocabCount = vocabCounts[dayKey] || 0;
          
          return (
            <div key={dayKey} className="relative">
              <button
                onClick={() => onDateSelect(day)}
                className={getDayClasses(day)}
              >
                {format(day, 'd')}
              </button>
              
              {vocabCount > 0 && (
                <div className="absolute -top-1 -right-1 bg-green-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                  {vocabCount > 9 ? '9+' : vocabCount}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Calendar;