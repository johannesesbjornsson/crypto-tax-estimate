import React, { useState, useEffect } from 'react';
import 'react-datepicker/dist/react-datepicker.css';

export default function UserProfile({ }) {
  const [user, setUser] = useState({ 
    name: '', 
    email: '', 
    currency: '',
    taxStartDate: '',
    taxEndDate: ''
  });

  useEffect(() => {
    async function fetchUserData() {
      try {
        const response = await fetch('/v1/user'); // Adjust URL accordingly
        if (!response.ok) throw new Error('Network response error');
        const data = await response.json();
        setUser(data);
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    }

    fetchUserData();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser(prev => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    try {
      const response = await fetch('/v1/user', {
        method: 'PUT', // or POST depending on your backend
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(user),
      });

      if (!response.ok) throw new Error('Failed to save user data');
      alert('Settings saved successfully');
    } catch (error) {
      console.error('Save failed:', error);
      alert('Error saving settings');
    }
  };

  return (
    <div className="settings">
      <div className="setting-row">
        <span className="setting-label">Email:</span>
        <span className="setting-value">{user.email}</span>
      </div>
      <div className="setting-row">
        <span className="setting-label">Name:</span>
        <input
          className="setting-value"
          type="text"
          name="name"
          value={user.name}
          onChange={handleChange}
        />
      </div>
      <div className="setting-row">
        <span className="setting-label">Tax Start Date:</span>
        
        <select className="setting-value" name="taxStartDay" value={user.taxStartDay} onChange={handleChange}>
          {Array.from({ length: 31 }, (_, i) => (
            <option key={i+1} value={i+1}>{i+1}</option>
          ))}
        </select>
        <span className="date-separator">/</span>
        <select className="setting-value" name="taxStartMonth" value={user.taxStartMonth} onChange={handleChange}>
          {Array.from({ length: 12 }, (_, i) => (
            <option key={i+1} value={i+1}>{i+1}</option>
          ))}
        </select>
      </div>
      <div className="setting-row">
        <span className="setting-label">Tax End Date:</span>
        <select className="setting-value" name="taxEndDay" value={user.taxEndDay} onChange={handleChange}>
          {Array.from({ length: 31 }, (_, i) => (
            <option key={i+1} value={i+1}>{i+1}</option>
          ))}
        </select>
        <span className="date-separator">/</span>
        <select className="setting-value" name="taxEndMonth" value={user.taxEndMonth} onChange={handleChange}>
          {Array.from({ length: 12 }, (_, i) => (
            <option key={i+1} value={i+1}>{i+1}</option>
          ))}
        </select>
      </div>
      <div className="setting-row">
        <span className="setting-label">Currency:</span>
        <select
          className="setting-value"
          value={user.currency}
          onChange={handleChange}
        >
          <option value="USD">USD</option>
          <option value="EUR">EUR</option>
          <option value="GBP">GBP</option>
          <option value="JPY">JPY</option>
        </select>
      </div>
      <div style={{ textAlign: 'right', marginTop: '20px' }}>
        <button onClick={handleSave}>Save</button>
      </div>
    </div>
  );
}
