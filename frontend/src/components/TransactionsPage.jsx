import React, { useState, useEffect } from 'react';
import 'react-datepicker/dist/react-datepicker.css';

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);

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

  return (
    <div className="transactions">
      <div className="tab-actions">
        <button onClick={() => alert('Add Transaction Clicked')}>Add Transaction</button>
        <button onClick={() => alert('Upload CSV Clicked')}>Upload CSV</button>
      </div>
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