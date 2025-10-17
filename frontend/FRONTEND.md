# Frontend Documentation

## Overview

The Vocabulary Tracker frontend is built with React 18, Vite, and TailwindCSS, providing a modern and responsive user interface for vocabulary management.

## Tech Stack

- **React 18**: Latest React with hooks and modern patterns
- **Vite**: Fast build tool and development server
- **TailwindCSS**: Utility-first CSS framework
- **React Router**: Client-side routing
- **Axios**: HTTP client for API calls
- **date-fns**: Date manipulation library
- **Lucide React**: Icon library

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Calendar.jsx     # Interactive calendar component
│   ├── VocabForm.jsx    # Form for adding/editing vocabulary
│   ├── VocabList.jsx    # List of vocabulary entries
│   └── ProtectedRoute.jsx # Route protection component
├── contexts/            # React contexts
│   └── AuthContext.jsx  # Authentication context
├── pages/               # Page components
│   ├── Login.jsx        # Login page
│   ├── Register.jsx     # Registration page
│   └── Dashboard.jsx    # Main dashboard
├── services/            # API services
│   └── api.js          # API client configuration
├── App.jsx             # Main app component
├── main.jsx            # App entry point
└── index.css           # Global styles with Tailwind
```

## Components

### Calendar Component

Interactive calendar for date selection and vocabulary visualization.

**Props:**
- `selectedDate`: Currently selected date
- `onDateSelect`: Function called when a date is clicked
- `vocabCounts`: Object with date keys and vocabulary counts

**Features:**
- Month navigation with arrow buttons
- Visual indicators for days with vocabulary entries
- Highlight today and selected date
- Responsive grid layout

### VocabForm Component

Modal form for creating and editing vocabulary entries.

**Props:**
- `vocab`: Existing vocabulary object (for editing)
- `onSave`: Function called when form is submitted
- `onCancel`: Function called when form is cancelled
- `selectedDate`: Date for new vocabulary entries

**Features:**
- Form validation with error messages
- Status selection (mastered/review needed)
- Modal overlay with backdrop click handling
- Loading states during API calls

### VocabList Component

Display and manage vocabulary entries for a selected date.

**Props:**
- `vocabularies`: Array of vocabulary objects
- `onEdit`: Function called when edit button is clicked
- `onDelete`: Function called when delete button is clicked
- `onUpdateStatus`: Function called when status is toggled
- `loading`: Boolean for loading state
- `searchTerm`: Current search term
- `onSearchChange`: Function called when search input changes

**Features:**
- Search functionality across words and meanings
- Status toggle buttons with visual indicators
- Edit and delete actions for each entry
- Empty state when no vocabulary found
- Loading skeleton during API calls