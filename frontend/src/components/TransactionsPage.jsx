import React, { useState } from 'react';

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState([
    { id: 1, date: '2024-05-01', description: 'Coffee', amount: -4.5 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
    { id: 2, date: '2024-05-02', description: 'Salary', amount: 2500 },
  ]);

  const handleRightClick = (transaction, x, y) => {
    alert(`Right-clicked on: ${transaction.description} at (${x}, ${y})`);
  // Or show a custom context menu here
  };

  return (
    <div className="transactions">
      <div className="tab-actions">
        <button onClick={() => alert('Add Transaction Clicked')}>Add Transaction</button>
        <button onClick={() => alert('Upload CSV Clicked')}>Upload CSV</button>
      </div>
      <div className="table-container">
      <table className="table">
        <thead>
          <tr>
            <th>Date</th>
            <th>Description</th>
            <th>Amount</th>
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
              <td>${tx.amount.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
      </div>
    </div>
  );
}

const thStyle = {
  textAlign: 'left',
  padding: '8px',
  borderBottom: '2px solid #ccc',
};

const tdStyle = {
  padding: '8px',
  borderBottom: '1px solid #eee',
};