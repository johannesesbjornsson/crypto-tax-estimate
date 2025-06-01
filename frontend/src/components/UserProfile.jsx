import React, { useState, useEffect } from 'react';
import Form from './common/Form';

export default function UserProfile() {
  const [successMessage, setSuccessMessage] = useState('');
  const [isError, setIsError] = useState(false);
  const [user, setUser] = useState({
    name: '',
    email: '',
    currency: '',
    taxStartDay: 1,
    taxStartMonth: 1,
    taxEndDay: 1,
    taxEndMonth: 1,
  });

  useEffect(() => {
    async function fetchUserData() {
      try {
        const response = await fetch('/v1/user');
        if (!response.ok) throw new Error('Network response error');
        const data = await response.json();
        setUser({
          ...data,
          taxStartDay: data.taxStartDay || 1,
          taxStartMonth: data.taxStartMonth || 1,
          taxEndDay: data.taxEndDay || 1,
          taxEndMonth: data.taxEndMonth || 1,
        });
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    }

    fetchUserData();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser((prev) => ({
      ...prev,
      [name]: name.includes('Day') || name.includes('Month') ? parseInt(value) : value,
    }));
  };

  const handleSave = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('/v1/user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user),
      });

      if (!response.ok) throw new Error('Failed to save user data');
      setSuccessMessage('✅ Settings saved');
      setIsError(false);
      setTimeout(() => setSuccessMessage(''), 3000);
    } catch (error) {
      console.error('Save failed:', error);
      setSuccessMessage('❌ Failed to save settings');
      setIsError(true);
    }
  };

  const days = Array.from({ length: 31 }, (_, i) => i + 1);
  const months = Array.from({ length: 12 }, (_, i) => i + 1);

  const userProfileFields = [
    { name: 'email', label: 'Email', type: 'display' },
    { name: 'name', label: 'Name', type: 'text' },
    {
      type: 'group',
      label: 'Tax Start Date',
      fields: [
        { name: 'taxStartDay', type: 'select', options: days },
        { name: 'taxStartMonth', type: 'select', options: months },
      ],
    },
    {
      type: 'group',
      label: 'Tax End Date',
      fields: [
        { name: 'taxEndDay', type: 'select', options: days },
        { name: 'taxEndMonth', type: 'select', options: months },
      ],
    },
    {
      name: 'currency',
      label: 'Currency',
      type: 'select',
      options: ['USD', 'EUR', 'GBP', 'JPY'],
    },
  ];

  return (
    <div className="settings">
      <Form
        fields={userProfileFields}
        formData={user}
        handleFormChange={handleChange}
        handleSubmit={handleSave}
        formErrorMessage={isError ? successMessage : ''}
        successMessage={!isError ? successMessage : ''}
        setShowForm={null}
      />
    </div>
  );
}