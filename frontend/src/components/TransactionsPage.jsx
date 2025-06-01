import React, { useState, useEffect } from 'react';
import 'react-datepicker/dist/react-datepicker.css';

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState([]);
  const [fileUploads, setFileUploads] = useState([]);
  const [successMessage, setSuccessMessage] = useState('');
  const [formErrorMessage, setFormErrorMessage] = useState('');
  const [isError, setIsError] = useState(false);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [showCSVForm, setShowCSVForm] = useState(false);
  const [showUploads, setShowUploads] = useState(false);
  const [page, setPage] = useState(1);
  const limit = 100;
  const [hasMore, setHasMore] = useState(false);

  const [formData, setFormData] = useState({
    date: '',
    description: '',
    type: 'Buy',
    amount: '',
    price: '',
    asset: '',
  });
  const [csvFile, setCsvFile] = useState(null);
  const [csvDescription, setCsvDescription] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const offset = (page - 1) * limit;

        const txRes = await fetch(`/v1/transactions?limit=${limit}&offset=${offset}`);
        const uploadRes = await fetch('/v1/transactions/upload');

        if (!txRes.ok || !uploadRes.ok) throw new Error('Failed to fetch data');

        const txData = await txRes.json();
        const fileData = await uploadRes.json();

        setTransactions(Array.isArray(txData) ? txData : []);
        setFileUploads(Array.isArray(fileData) ? fileData : []);
        setHasMore(txData.length === limit);
      } catch (error) {
        console.error('Error:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [page]);

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
          price: parseFloat(formData.price),
          date: new Date(formData.date).toISOString(),
        }),
      });
      if (!res.ok) throw new Error('Failed to add transaction');
      const updated = await res.json();
      setTransactions(prev => [updated, ...prev]);
      setFormData({ date: '', description: '', price: '', type: 'Buy', amount: '', asset: '' });
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

  const handleCSVUpload = async (e) => {
    e.preventDefault();
    setFormErrorMessage('');
    if (!csvFile) {
      setFormErrorMessage('❌ Please choose a file');
      return;
    }

    try {
      const form = new FormData();
      form.append('file', csvFile);
      form.append('description', csvDescription);

      const res = await fetch('/v1/transactions/upload', {
        method: 'POST',
        body: form,
      });

      if (!res.ok) throw new Error('Failed to upload CSV');

      setCsvFile(null);
      setCsvDescription('');
      setShowCSVForm(false);
      setSuccessMessage('✅ CSV uploaded successfully');
      setIsError(false);
      setTimeout(() => setSuccessMessage(''), 3000);

      setPage(1); // reset to first page after upload
    } catch (error) {
      console.error('CSV upload failed:', error);
      setFormErrorMessage('❌ Failed to upload CSV');
      setIsError(true);
    }
  };

  const handleAddTransactionClick = () => {
    setShowForm(true);
    setShowCSVForm(false);
  };

  const handleCSVUploadClick = () => {
    setShowCSVForm(true);
    setShowForm(false);
  };

  return (
    <div className="transactions">
      <div className="tab-actions">
        <button onClick={handleAddTransactionClick}>Add Transaction</button>
        <button onClick={handleCSVUploadClick}>Upload CSV</button>
      </div>

      {successMessage && !showForm && !showCSVForm && (
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
              <input className="setting-value" name="description" value={formData.description} onChange={handleFormChange} />
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
              <input className="setting-value" name="amount" type="number" step="0.01" value={formData.amount} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Price:</span>
              <input className="setting-value" name="price" type="number" step="0.01" value={formData.price} onChange={handleFormChange} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Asset:</span>
              <input className="setting-value" name="asset" value={formData.asset} onChange={handleFormChange} required />
            </div>
            <div className="save-row">
              {formErrorMessage && <span className="save-message error">{formErrorMessage}</span>}
              <button type="submit">Save</button>
              <button type="button" onClick={() => setShowForm(false)}>Cancel</button>
            </div>
          </form>
        </div>
      )}

      {showCSVForm && (
        <div className="settings">
          <form onSubmit={handleCSVUpload}>
            <div className="setting-row">
              <span className="setting-label">File:</span>
              <input className="setting-value" type="file" accept=".csv" onChange={(e) => setCsvFile(e.target.files[0])} required />
            </div>
            <div className="setting-row">
              <span className="setting-label">Description:</span>
              <input className="setting-value" type="text" value={csvDescription} onChange={(e) => setCsvDescription(e.target.value)} />
            </div>
            <div className="save-row">
              {formErrorMessage && <span className="save-message error">{formErrorMessage}</span>}
              <button type="submit">Upload</button>
              <button type="button" onClick={() => setShowCSVForm(false)}>Cancel</button>
            </div>
          </form>
        </div>
      )}

      {!loading && fileUploads.length > 0 && (
        <div className="tab-actions">
          <button onClick={() => setShowUploads(prev => !prev)}>
            {showUploads ? 'Hide Uploaded Files' : 'Show Uploaded Files'}
          </button>
        </div>
      )}

      {showUploads && !loading && fileUploads.length > 0 && (
        <div className="table-container" style={{ marginBottom: '4rem' }}>
          <h3>Uploaded Files</h3>
          <table className="table">
            <thead>
              <tr>
                <th>Filename</th>
                <th>Uploaded At</th>
              </tr>
            </thead>
            <tbody>
              {fileUploads.map(file => (
                <tr key={file.id}>
                  <td>{file.name}</td>
                  <td>{new Date(file.created_at).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="table-container">
        <h3>Transactions</h3>
        {loading ? (
          <p>Loading transactions...</p>
        ) : (
          <>
            <table className="table">
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Description</th>
                  <th>Type</th>
                  <th>Amount</th>
                  <th>Price</th>
                  <th>Asset</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map(tx => (
                  <tr key={tx.id} onContextMenu={(e) => {
                    e.preventDefault();
                    handleRightClick(tx, e.clientX, e.clientY);
                  }}>
                    <td>{tx.date}</td>
                    <td>{tx.description}</td>
                    <td>{tx.type}</td>
                    <td>{tx.amount.toFixed(10)}</td>
                    <td>{tx.price.toFixed(5)}</td>
                    <td>{tx.asset}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div style={{ marginTop: '1rem' }}>
              <button onClick={() => setPage(p => Math.max(p - 1, 1))} disabled={page === 1}>Previous</button>
              <span style={{ margin: '0 1rem' }}>Page {page}</span>
              <button onClick={() => setPage(p => p + 1)} disabled={!hasMore}>Next</button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}