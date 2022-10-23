import React from 'react';
import { Theme, makeStyles } from '@material-ui/core/styles';
import { Box, Typography } from '@material-ui/core';
import { drawerWidth } from '../../../../common/shared/dimensions';

const useStyles = makeStyles((theme: Theme) => ({
  logoWrapper: {
    width: `calc(${drawerWidth}px - 48px)`, // 240 - 24*2
    textAlign: 'center',
    display: 'flex',
    alignItems: 'center',
    '& svg': {
      marginRight: 8
    }
  },
  landscapeIcon: {
    fill: "url(#landscapeGradient)",
    height: theme.spacing(4),
    width: theme.spacing(4)
  }
}));

const MainLogo: React.FC = () => {
  const classes = useStyles();

  return (
    <Box className={classes.logoWrapper}>
      <Typography variant="h5">Crypto Insight</Typography>
    </Box>
  )
}

export default MainLogo;
