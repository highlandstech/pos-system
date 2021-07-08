import React, { useState } from 'react';
import styled from 'styled-components';
import {
    AppBar,
    Button,
    CircularProgress,
    Container,
    Paper,
    Tab,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Tabs,
    Typography,
} from '@material-ui/core';
import { RichOrder } from '../RichOrder';
import { EarningsCard } from '../EarningsCard';
import {
    useGetOrdersList,
    useGetEarningsPerDay,
    useGetTotalEarnings,
} from '../../hooks';

const GridRowBox = styled.div`
    display: grid;
    grid-template-columns: 25% 25% 25% 25%;
    margin-bottom: 15px;

    & > div:not(:last-child) {
        margin-right: 5px;
    }
`;

export interface IMainProps {};

interface IMainState {
    page: number;
    tab: number;
}

/**
 * Renders the main page.
 * @param props The props.
 * @returns The Main page component.
 */
export const Main = (props: IMainProps) => {
    const [state, setState] = useState<IMainState>({
        page: 1,
        tab: 0,
    });
    const { loading, error, orders, fetchOrders } = useGetOrdersList(state.page);
    const earningsPerDay = useGetEarningsPerDay();
    const earnings = useGetTotalEarnings();

    const exportTotals = () => {
        const link = document.createElement('a');
        link.download = '';
        link.href = '/api/orders/totals/export';
        link.click();
    };

    return (
        <div>
            <AppBar>
                <Tabs value={state.tab} onChange={(_, tab) => setState({ ...state, tab })}>
                    <Tab label="Orders" />
                    <Tab label="Users" />
                </Tabs>
            </AppBar>
            <br />
            <br />
            <br />
            {state.tab === 0 ? (
                <Container>
                    <Typography variant="h3" component="h3" gutterBottom>
                        Total Earnings: ${earnings.toFixed(2)}
                    </Typography>
                    <GridRowBox>
                        <EarningsCard day={0} amount={earningsPerDay[0]} />
                        <EarningsCard day={1} amount={earningsPerDay[1]} />
                        <EarningsCard day={2} amount={earningsPerDay[2]} />
                        <EarningsCard day={3} amount={earningsPerDay[3]} />
                    </GridRowBox>
                    <div style={{ marginBottom: '15px' }}>
                        <Button onClick={exportTotals} variant="contained" color="primary">
                            Export Sales YTD
                        </Button>
                    </div>
                    <div>
                        {error ? (
                            <Typography variant="h4" component="h4">
                                {error}
                            </Typography>
                        ) : null}

                        {loading ? (
                            <div style={{ textAlign: 'center' }}>
                                <CircularProgress
                                    variant="indeterminate"
                                    disableShrink
                                    style={{
                                        strokeLinecap: 'round',
                                        color: '#1a90ff',
                                        animationDuration: '550ms',
                                        marginTop: '20px',
                                    }}
                                    size={40}
                                    thickness={4}
                                />
                            </div>
                        ) : null}

                        <TableContainer component={Paper}>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>#</TableCell>
                                        <TableCell>Purchased Items</TableCell>
                                        <TableCell></TableCell>
                                        <TableCell>Total</TableCell>
                                        <TableCell>Date Placed</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {orders.map((order, i) => <RichOrder key={i} order={order} />)}
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </div>
                    <br />
                </Container>
            ) : null}

            {state.tab === 1 ? (
                <div>Hello</div>
            ) : null}
        </div>
    );
};
