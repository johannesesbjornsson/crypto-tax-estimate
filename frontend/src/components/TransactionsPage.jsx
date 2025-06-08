import React, { useState, useEffect } from 'react';
import 'react-datepicker/dist/react-datepicker.css';
import Form from './common/Form';
import Table from './common/Table';

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
  const [totalPages, setTotalPages] = useState(1);

  const [formData, setFormData] = useState({
    type: 'Buy',
    date: '',
    description: '',
    amount: '',
    price: '',
    quote_currency: 'USD',
    asset: '',
  });

  const [csvFile, setCsvFile] = useState(null);
  const [csvDescription, setCsvDescription] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const offset = (page - 1) * limit;

        const txRes = await fetch(`/v1/transactions?limit=${limit}&offset=${offset}&txType=trade`);
        const uploadRes = await fetch('/v1/transactions/upload');

        if (!txRes.ok || !uploadRes.ok) throw new Error('Failed to fetch data');

        const txData = await txRes.json();
        const fileData = await uploadRes.json();

        setTransactions(Array.isArray(txData.transactions) ? txData.transactions : []);
        setTotalPages(txData.totalPages || 1);
        setHasMore(page < txData.totalPages);
        setFileUploads(Array.isArray(fileData) ? fileData : []);
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
    const { name, value, files } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'file' ? files[0] : value,
    }));
  };
  const handleCSVFormChange = (e) => {
    const { name, value, files } = e.target;

    if (name === 'file') {
      setCsvFile(files[0]);
    } else {
      setCsvDescription(value);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFormErrorMessage('');

    try {
      const { type, date, amount, price, quote_currency, description, asset } = formData;
      const isTrade = type === 'Buy' || type === 'Sell';

      const payload = {
        type,
        date: new Date(date).toISOString(),
        amount: parseFloat(amount),
        asset,
        description,
      };

      let url = '/v1/transactions';
      if (isTrade) {
        const query = new URLSearchParams({
          price: price.toString(),
          quote_currency,
        }).toString();
        url += `?${query}`;
      }

      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (!res.ok) throw new Error('Failed to add transaction');
      const updated = await res.json();
      setTransactions(prev => [updated, ...prev]);

      setFormData({
        date: '',
        description: '',
        price: '',
        quote_currency: 'USD',
        type: 'Buy',
        amount: '',
        asset: ''
      });

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

      setPage(1);
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

  const isTrade = formData.type === 'Buy' || formData.type === 'Sell';

  const baseFields = [
    { name: 'type', label: 'Type', type: 'select', options: ['Income', 'Buy', 'Sell', 'Lost'] },
    { name: 'date', label: 'Date', type: 'date', required: true },
    { name: 'description', label: 'Description', type: 'text' },
    { name: 'amount', label: 'Amount', type: 'number', step: '0.01', required: true },
    { name: 'asset', label: 'Asset', type: 'text', required: true },
  ];

  const tradeFields = [
    { name: 'price', label: 'Price', type: 'number', step: '0.01', required: true },
    { name: 'quote_currency', label: 'Quote Currency', type: 'text', required: true },
  ];

  const fields = isTrade ? [...baseFields, ...tradeFields] : baseFields;

  const csvFormData = {
    file: csvFile,
    description: csvDescription
  };

  const csvFormFields = [
    { name: 'file', label: 'File', type: 'file', accept: '.csv', required: true },
    { name: 'description', label: 'Description', type: 'text' }
  ];

  const tableColumns = [
    { key: 'date', label: 'Date' },
    { key: 'description', label: 'Description' },
    { key: 'type', label: 'Type' },
    { key: 'amount', label: 'Amount', render: val => parseFloat(val).toFixed(10) },
    { key: 'price', label: 'Price', render: val => parseFloat(val).toFixed(5) },
    { key: 'asset', label: 'Asset' },
  ];

  return (
    <div className="transactions">
      <div className="tab-actions">
        <button onClick={handleAddTransactionClick}>Add Transaction</button>
        <button onClick={handleCSVUploadClick}>Upload CSV</button>
      </div>

      {successMessage && !showForm && !showCSVForm && (
        <div className={`save-message ${isError ? 'error' : 'success'}`}>
          {successMessage}
        </div>
      )}

      {showForm && (
        <div className="settings">
          <Form
            fields={fields}
            formData={formData}
            handleFormChange={handleFormChange}
            handleSubmit={handleSubmit}
            formErrorMessage={formErrorMessage}
            setShowForm={setShowForm}
          />
        </div>
      )}

      {showCSVForm && (
        <div className="settings">
          <Form
            fields={csvFormFields}
            formData={csvFormData}
            handleFormChange={handleCSVFormChange}
            handleSubmit={handleCSVUpload}
            formErrorMessage={formErrorMessage}
            setShowForm={setShowCSVForm}
          />
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
        <div className="table-container">
          <h3>Uploaded Files</h3>
          <Table
            columns={[
              { key: 'name', label: 'Filename' },
              { key: 'created_at', label: 'Uploaded At', render: val => new Date(val).toLocaleString() },
            ]}
            data={fileUploads}
          />
        </div>
      )}

      <div className="table-container">
        <h3>Transactions</h3>
        {loading ? (
          <p>Loading transactions...</p>
        ) : (
          <>
            <Table
              columns={tableColumns}
              data={transactions}
              onRowContextMenu={handleRightClick}
            />
            <div style={{ marginTop: '1rem' }}>
              <button onClick={() => setPage(p => Math.max(p - 1, 1))} disabled={page === 1}>
                Previous
              </button>
              <span style={{ margin: '0 1rem' }}>
                Page {page} of {totalPages}
              </span>
              <button onClick={() => setPage(p => p + 1)} disabled={page >= totalPages}>
                Next
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}