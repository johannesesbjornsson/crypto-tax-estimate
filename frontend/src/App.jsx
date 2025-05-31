import { useState, useEffect } from 'react'
import './App.css'
import TransactionsPage from './components/TransactionsPage';

import Button from './components/Button';
import Tabs from './components/Tabs';
import UserProfile from './components/UserProfile';


function App() {
  const [activeTab, setActiveTab] = useState(() => {
    return localStorage.getItem('activeTab') || 'Overview';
  });

  useEffect(() => {
    localStorage.setItem('activeTab', activeTab);
  }, [activeTab]);
  const tabs = {
    Overview: <div className="">Overview content</div>,
        Transactions: <TransactionsPage />,
        Profile: (
      <UserProfile/>
    ),
  };

  return (
    <div className="">
      <Tabs tabs={tabs} activeTab={activeTab} onTabChange={setActiveTab} />
    </div>
  );
}

export default App;
