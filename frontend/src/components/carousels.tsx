import { useState } from 'react';
import { Carousel } from 'react-bootstrap';

const Carousels = () => {
  const [index, setIndex] = useState(0);

  const handleSelect = (selectedIndex: number) => {
    setIndex(selectedIndex);
  };

  return (
    <Carousel activeIndex={index} onSelect={handleSelect}>
      <Carousel.Item className='carsouel__item'>
        <img
          className='d-block w-full '
          src='/images/1.jpg'
          alt='First slide'
        />
      </Carousel.Item>
      <Carousel.Item className='carsouel__item'>
        <img
          className='d-block w-full '
          src='/images/2.jpg'
          alt='Second slide'
        />
      </Carousel.Item>
      <Carousel.Item className='carsouel__item'>
        <img
          className='d-block w-full '
          src='/images/3.png'
          alt='Third slide'
        />
      </Carousel.Item>
    </Carousel>
  );
};

export default Carousels;
