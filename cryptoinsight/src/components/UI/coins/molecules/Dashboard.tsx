import React, { useState, useEffect } from "react";
import { Theme, makeStyles } from "@material-ui/core/styles";
import { Grid } from "@material-ui/core";
import SocialCard from "../atoms/SocialCard";
import { useAppSelector } from "../../../../app/hooks";
import { selectCoinDetails } from "../../../../features/coinDetailsSlice";
import { Notifications, NotificationsActive } from "@material-ui/icons";

import jwt from "jsonwebtoken";

const useStyles = makeStyles((theme: Theme) => ({
  container: {
    minHeight: 600,
  },
}));

const Dashboard: React.FC = () => {
  const classes = useStyles();

  const loadIframeURL = (dashboard: number) => {
    var METABASE_SITE_URL = "http://ec2-13-230-149-127.ap-northeast-1.compute.amazonaws.com:3000";
    var METABASE_SECRET_KEY = "c68d6a4cb0619b8314d521f39b98fd9494a3a63c2c3a8362b2491d255c7bc333";

    var payload = {
      resource: { dashboard: dashboard },
      params: {},
      exp: Math.round(Date.now() / 1000) + 10 * 60, // 10 minute expiration
    };
    var token = jwt.sign(payload, METABASE_SECRET_KEY);
    var iframeUrl = METABASE_SITE_URL + "/embed/dashboard/" + token + "#bordered=false&titled=false";
    return iframeUrl;
  };

  const [iframeUrl1, setIframeUrl1] = useState<string>("");
  const [iframeUrl2, setIframeUrl2] = useState<string>("");
  const [iframeUrl3, setIframeUrl3] = useState<string>("");
  const [height, setHeight] = useState<number>(410);

  const [active1, setActive1] = useState<boolean>(false);
  const [active2, setActive2] = useState<boolean>(false);
  const [active3, setActive3] = useState<boolean>(false);

  const coinDetails = useAppSelector(selectCoinDetails);
  console.log(coinDetails);

  useEffect(() => {
    if (coinDetails.value?.id === "bitcoin") {
      setHeight(420);
      setIframeUrl1(loadIframeURL(10));
      setIframeUrl2(loadIframeURL(9));
      setIframeUrl3(loadIframeURL(7));
    } else if (coinDetails.value?.id === "ethereum") {
      setHeight(380);
      setIframeUrl1(loadIframeURL(5));
      setIframeUrl2(loadIframeURL(8));
      setIframeUrl3(loadIframeURL(11));
    }
  }, [coinDetails]);

  return (
    <Grid container spacing={3} className={classes.container}>
      <Grid item xs={12}>
        <SocialCard
          title={coinDetails.value?.id === "bitcoin" ? " " : " "}
          icon={active1 ? <NotificationsActive /> : <Notifications />}
          iconColor={active1 ? "#FF4500" : "#333333"}
          onClick={() => {
            setActive1(!active1);
          }}
          showAction={true}
          link={iframeUrl1}
        >
          <iframe
            title="Insight"
            src={iframeUrl1}
            width="100%"
            height={height}
            allowTransparency={true}
            allowFullScreen={true}
            style={{ border: 0 }}
          />
        </SocialCard>
      </Grid>
      <Grid item xs={12}>
        <SocialCard
          title={coinDetails.value?.id === "bitcoin" ? " " : " "}
          icon={active2 ? <NotificationsActive /> : <Notifications />}
          iconColor={active2 ? "#FF4500" : "#333333"}
          onClick={() => {
            setActive2(!active2);
          }}
          link={iframeUrl2}
        >
          <iframe
            title="Insight"
            src={iframeUrl2}
            width="100%"
            height={height}
            allowTransparency={true}
            allowFullScreen={true}
            style={{ border: 0 }}
          />
        </SocialCard>
      </Grid>
      <Grid item xs={12}>
        <SocialCard
          title={coinDetails.value?.id === "bitcoin" ? " " : " "}
          icon={active3 ? <NotificationsActive /> : <Notifications />}
          iconColor={active3 ? "#FF4500" : "#333333"}
          onClick={() => {
            setActive3(!active3);
          }}
          link={iframeUrl3}
        >
          <iframe
            title="Insight"
            src={iframeUrl3}
            width="100%"
            height={height}
            allowTransparency={true}
            allowFullScreen={true}
            style={{ border: 0 }}
          />
        </SocialCard>
      </Grid>
    </Grid>
  );
};

export default Dashboard;
