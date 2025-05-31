import React, { useState } from 'react';

export default function Tabs({ tabs }) {
  const [activeTab, setActiveTab] = useState(() => {
    return localStorage.getItem('activeTab') || Object.keys(tabs)[0];
  });

  return (
    <>
      <div className="tab-bar">
        {Object.keys(tabs).map((key) => (
          <button
            key={key}
            onClick={() => {
              localStorage.setItem('activeTab', key);
              setActiveTab(key);
            }}
            className={`px-4 py-2${key === activeTab ? ' active' : ''}`}
          >
            {key}
          </button>
        ))}
      </div>

      <div className="">
        {tabs[activeTab]}
      </div>
    </>
  );
}