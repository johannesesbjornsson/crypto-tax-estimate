import React, { useState } from 'react';

export default function Tabs({ tabs }) {
  const [activeTab, setActiveTab] = useState(Object.keys(tabs)[0]);

  return (
    <>
      {/* Fixed tab bar at the top of the screen */}
      <div className="tab-bar">
        {Object.keys(tabs).map((key) => (
          <button
            key={key}
            onClick={() => setActiveTab(key)}
            className={`px-4 py-2${key === activeTab ? ' active' : ''}`}
          >
            {key}
          </button>
        ))}
      </div>

      {/* Add enough padding to push content below the fixed tabs */}
      <div className="">
        {tabs[activeTab]}
      </div>
    </>
  );
}