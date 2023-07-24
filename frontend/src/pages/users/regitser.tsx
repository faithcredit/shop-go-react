import React from 'react';
import FormContainer from '../../components/UI/form-container';
import { Button, Form } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import { FormProvider,useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import publicAxios from '../../utils/public-axios';
import {toast} from 'react-toastify';
import { setError } from '../../utils/error';

import { object, string, TypeOf } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

type FormValues = {
  name: string;
  email: string;
  password: string;
  passwordConfirm: string;
};

const validationSchema = object({
  name: string().min(1, 'Full name is required').max(100),
  email: string()
    .min(1, 'Email address is required')
    .email('Email Address is invalid'),
  password: string()
    .min(1, 'Password is required')
    .min(8, 'Password must be more than 8 characters')
    .max(32, 'Password must be less than 32 characters'),
  passwordConfirm: string().min(1, 'Please confirm your password'),
}).refine((data) => data.password === data.passwordConfirm, {
  path: ['passwordConfirm'],
  message: 'Passwords do not match',
});

export type RegisterInput = TypeOf<typeof validationSchema>;

const Register = () => {
  const navigate = useNavigate();

  const methods = useForm<RegisterInput>({
    resolver: zodResolver(validationSchema),
  });

  const {
    reset,
    register,
    handleSubmit,
    formState: { errors },
  } = methods

  const onSubmit = (data: FormValues) => {
    publicAxios
      .post('/auth/register', data)
      .then((res) => {
        toast.success('you have been registred , please login');
        navigate('/login');
      })
      .catch((err) => toast.error(setError(err)));
  };

  return (
    <FormContainer
      meta='register for free'
      image='/images/p2.jpg'
      title='Register For Free'
    >
      <FormProvider {...methods}>
        <Form onSubmit={handleSubmit(onSubmit)} noValidate>
          <Form.Group controlId='name'>
            <Form.Label>Username</Form.Label>
            <Form.Control
              placeholder='Enter name'
              {...register('name')}
              className={errors.name?.message && 'is-invalid'}
            />
            <p className='invalid-feedback'>{errors.name?.message}</p>
          </Form.Group>
          <Form.Group controlId='email'>
            <Form.Label>Email</Form.Label>

            <Form.Control
              type='email'
              placeholder='Enter email'
              {...register('email')}
              className={errors.email?.message && 'is-invalid'}
            />
            <p className='invalid-feedback'>{errors.email?.message}</p>
          </Form.Group>
          <Form.Group controlId='password'>
            <Form.Label>Mot de Passe </Form.Label>

            <Form.Control
              type='password'
              placeholder='*******'
              {...register('password')}
              className={errors.password?.message && 'is-invalid'}
            />
            <p className='invalid-feedback'>{errors.password?.message}</p>
          </Form.Group>
          <Form.Group controlId='passwordConfirm'>
            <Form.Label>Confirm Password </Form.Label>

            <Form.Control
              type='password'
              placeholder='*******'
              {...register('passwordConfirm')}
              className={errors.passwordConfirm?.message && 'is-invalid'}
            />
            <p className='invalid-feedback'>{errors.passwordConfirm?.message}</p>
            <Link to='/login' className='float-end me-2 mt-1'>
              Already have an Account ? Login
            </Link>
          </Form.Group>

          <Button
            style={{ backgroundColor: '#0071dc', color: '#fff' }}
            variant='outline-none'
            type='submit'
            className='mt-4 w-full'
          >
            Register
          </Button>
        </Form>
      </FormProvider>
    </FormContainer>
  );
};

export default Register;
