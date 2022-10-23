import React from 'react';
import {
  AccountBalanceRounded,
  DataUsageRounded,
  TollRounded,
  TrendingUpRounded
} from '@material-ui/icons';
import { BrowserRouter } from 'react-router-dom';
import PageLayout from '../components/templates/PageLayout';
import Overview from './analytics/Overview';
import Exchanges from './explore/Exchanges';
import Trends from './analytics/Trends';
import Coins from './explore/Coins';
import CoinDetails from './explore/CoinDetails';
import { RootModule } from '../models';

const Main: React.FC = () => {
  const rootModule: RootModule[] = [
    {
      moduleName: 'Analytics',
      pages: [
        {
          label: 'Overview',
          path: '/',
          icon: <DataUsageRounded />,
          page: <Overview />,
          index: 0
        },
        {
          label: 'Trends',
          path: '/trends',
          icon: <TrendingUpRounded />,
          page: <Trends />,
          index: 1
        }
      ]
    },
    {
      moduleName: 'Explore',
      pages: [
        {
          label: 'Coins',
          path: '/coins',
          icon: <TollRounded />,
          page: <Coins />,
          index: 2,
          subpage: {
            path: 'coinId',
            page: <CoinDetails />
          }
        },
        {
          label: 'Exchanges',
          path: '/exchanges',
          icon: <AccountBalanceRounded />,
          page: <Exchanges />,
          index: 3

        },
      ]
    }
  ];

  return (
    <BrowserRouter>
      <PageLayout rootModule={rootModule} />
    </BrowserRouter>
  );
}

export default Main;