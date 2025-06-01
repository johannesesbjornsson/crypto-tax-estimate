import React from 'react';

export default function Table({ columns, data, onRowContextMenu }) {
  if (!data || data.length === 0) {
    return <p>No data available.</p>;
  }

  return (
    <table className="table">
      <thead>
        <tr>
          {columns.map(col => (
            <th key={col.key}>{col.label}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((row, rowIndex) => (
          <tr
            key={row.id || rowIndex}
            onContextMenu={onRowContextMenu ? (e) => {
              e.preventDefault();
              onRowContextMenu(row, e.clientX, e.clientY);
            } : undefined}
          >
            {columns.map(col => (
              <td key={col.key}>
                {col.render ? col.render(row[col.key], row) : row[col.key]}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}