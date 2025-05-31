import React, { useState, useEffect } from 'react';
import 'react-datepicker/dist/react-datepicker.css';

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState([]);
  const [successMessage, setSuccessMessage] = useState('');
  const [formErrorMessage, setFormErrorMessage] = useState('');
  const [isError, setIsError] = useState(false);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    date: '',
    description: '',
    venue: '',
    type: 'Buy',
    amount: '',
    asset: '',
    source: '',
  });

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const response = await fetch(`/v1/transactions`);
        if (!response.ok) throw new Error('Failed to fetch transactions');
        const data = await response.json();
        setTransactions(data);
      } catch (error) {
        console.error('Error:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchTransactions();
  }, []);

  const handleRightClick = (transaction, x, y) => {
    alert(`Right-clicked on: ${transaction.description} at (${x}, ${y})`);
  };

  const handleFormChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFormErrorMessage('');
    try {
      const res = await fetch('/v1/transactions', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ...formData,
          amount: parseFloat(formData.amount),
          date: new Date(formData.date).toISOString(),
        }),
      });
      if (!res.ok) throw new Error('Failed to add transaction');
      const updated = await res.json();
      setTransactions(prev => [...prev, updated]);
      setFormData({ date: '', description: '', venue: '', type: 'Buy', amount: '', asset: '', source: '' });
      setShowForm(false);
      setSuccessMessage('✅ Transaction added');
      setIsError(false);
      setTimeout(() => setSuccessMessage(''), 3000);
    } catch (error) {
      console.error('Save failed:', error);
      setFormErrorMessage('❌ Failed to save transaction');
      setIsError(true);
    }
  };

  return (
    <div className="transactions">
      <div className="tab-actions">
        <button onClick={() => setShowForm(true)}>Add Transaction</button>
        <button onClick={() => alert('Upload CSV Clicked')}>Upload CSV</button>
      </div>

      {successMessage && !showForm && (
        <div className={`save-message ${isError ? 'error' : 'success'}`} style={{ marginBottom: '1rem' }}>
          {successMessage}
        </div>
      )}

      {showForm && (
        <div className="settings">
          <form onSubmit={handleSubmit}>
            <div className="setting-row">
              <span className="setting-label">Date:</span>
              <input className="setting-value" name="date" type="date" value={formData.date} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Description:</span>
              <input className="setting-value" name="description" placeholder="Description" value={formData.description} onChange={handleFormChange} />
            </div>
            <div className="setting-row">
              <span className="setting-label">Venue:</span>
              <input className="setting-value" name="venue" placeholder="Venue" value={formData.venue} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Type:</span>
              <select className="setting-value" name="type" value={formData.type} onChange={handleFormChange}>
                <option>Income</option>
                <option>Buy</option>
                <option>Sell</option>
                <option>Lost</option>
              </select>
            </div>
            <div className="setting-row">
              <span className="setting-label">Amount:</span>
              <input className="setting-value" name="amount" type="number" step="0.01" placeholder="Amount" value={formData.amount} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Asset:</span>
              <input className="setting-value" name="asset" placeholder="Asset" value={formData.asset} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Source:</span>
              <input className="setting-value" name="source" placeholder="Source" value={formData.source} onChange={handleFormChange} />
            </div>

<div className="save-row" style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
  {formErrorMessage && (
    <span className="save-message error" style={{ whiteSpace: 'nowrap' }}>
      {formErrorMessage}
    </span>
  )}
  <button type="submit">Save</button>
  <button type="button" onClick={() => setShowForm(false)}>Cancel</button>
</div>
          </form>
        </div>
      )}

      <div className="table-container">
        {loading ? (
          <p>Loading transactions...</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Venue</th>
                <th>Type</th>
                <th>Amount</th>
                <th>Asset</th>
                <th>Source</th>
              </tr>
            </thead>
            <tbody>
              {transactions.map(tx => (
                <tr
                  key={tx.id}
                  onContextMenu={(e) => {
                    e.preventDefault();
                    handleRightClick(tx, e.clientX, e.clientY);
                  }}
                >
                  <td>{tx.date}</td>
                  <td>{tx.description}</td>
                  <td>{tx.venue}</td>
                  <td>{tx.type}</td>
                  <td>${tx.amount.toFixed(2)}</td>
                  <td>{tx.asset}</td>
                  <td>{tx.source}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}