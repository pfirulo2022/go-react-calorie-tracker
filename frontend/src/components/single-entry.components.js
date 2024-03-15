
//import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import 'bootstrap/dist/css/bootstrap.css';
import { Button, Card, Row, Col } from 'react-bootstrap';



export const Entry = ({ entryData, setChangeIngredient, deleteSingleEntry, setChangeEntry }) => {
  return (
    <Card>
      <Row>
        <Col>Dish:{entryData !== undefined && entryData.dish}  </Col>
        <Col>Dish:{entryData !== undefined && entryData.ingredients}  </Col>
        <Col>Dish:{entryData !== undefined && entryData.calories}  </Col>
        <Col>Dish:{entryData !== undefined && entryData.fat}  </Col>
        <Col>
          <Button variant="primary" onClick={() => deleteSingleEntry(entryData._id)}>Delete Entry</Button>
        </Col>
        <Col>
          <Button variant="primary" onClick={() => (ChangeIngredient())}>Change ingredients </Button>
        </Col>
        <Col>
          <Button variant="primary" onClick={() => ChangeEntry()}>Change Entry</Button>
        </Col>
      </Row>

    </Card>
  )

  function ChangeIngredient() {
    setChangeIngredient({
"change":true,
"id":entryData._id
    })
  }
}



Entry.propTypes = {
  entryData: PropTypes.object,
  setChangeIngredient: PropTypes.func,
  deleteSingleEntry: PropTypes.func,
  setChangeEntry: PropTypes.func
}