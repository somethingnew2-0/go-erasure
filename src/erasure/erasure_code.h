/**********************************************************************
  Copyright(c) 2011-2014 Intel Corporation All rights reserved.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions 
  are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in
      the documentation and/or other materials provided with the
      distribution.
    * Neither the name of Intel Corporation nor the names of its
      contributors may be used to endorse or promote products derived
      from this software without specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
  "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
  LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
  A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
  OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
  SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
  LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
  DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
  THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
  OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
**********************************************************************/


#ifndef _ERASURE_CODE_H_
#define _ERASURE_CODE_H_

/**
 *  @file erasure_code.h
 *  @brief Interface to functions supporting erasure code encode and decode.
 *
 *  This file defines the interface to optimized functions used in erasure
 *  codes.  Encode and decode of erasures in GF(2^8) are made by calculating the
 *  dot product of the symbols (bytes in GF(2^8)) across a set of buffers and a
 *  set of coefficients.  Values for the coefficients are determined by the type
 *  of erasure code.  Using a general dot product means that any sequence of
 *  coefficients may be used including erasure codes based on random 
 *  coefficients.
 *  Multiple versions of dot product are supplied to calculate 1-6 output
 *  vectors in one pass.
 *  Base GF multiply and divide functions can be sped up by defining
 *  GF_LARGE_TABLES at the expense of memory size.
 *
 */

#include "gf_vect_mul.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Initialize tables for fast Erasure Code encode and decode.
 *
 * Generates the expanded tables needed for fast encode or decode for erasure
 * codes on blocks of data.  32bytes is generated for each input coefficient.
 *
 * @param k      The number of vector sources or rows in the generator matrix
 *               for coding.
 * @param rows   The number of output vectors to concurrently encode/decode.
 * @param a      Pointer to sets of arrays of input coefficients used to encode
 *               or decode data.
 * @param gftbls Pointer to start of space for concatenated output tables
 *               generated from input coefficients.  Must be of size 32*k*rows.
 * @returns none
 */

void ec_init_tables(int k, int rows, unsigned char* a, unsigned char* gftbls);

/**
 * @brief Generate or decode erasure codes on blocks of data, runs baseline version.
 *
 * Given a list of source data blocks, generate one or multiple blocks of
 * encoded data as specified by a matrix of GF(2^8) coefficients.  When given a
 * suitable set of coefficients, this function will perform the fast generation
 * or decoding of Reed-Solomon type erasure codes.
 *
 * @param len    Length of each block of data (vector) of source or dest data.
 * @param srcs   The number of vector sources or rows in the generator matrix
 * 		 for coding.
 * @param dests  The number of output vectors to concurrently encode/decode.
 * @param v      Pointer to array of input tables generated from coding
 * 		 coefficients in ec_init_tables(). Must be of size 32*k*rows
 * @param src    Pointer to source input buffer.
 * @param dest   Pointer to coded output buffers.
 * @returns none
 */

void ec_encode_data(int len, int srcs, int dests, unsigned char *v, unsigned char *src, unsigned char *dest);


/**
 * @brief GF(2^8) vector dot product
 * 
 * Does a GF(2^8) dot product across each byte of the input array and a constant
 * set of coefficients to produce each byte of the output. Can be used for
 * erasure coding encode and decode. Function requires pre-calculation of a
 * 32*vlen byte constant array based on the input coefficients.
 * 
 * @param len    Length of each vector in bytes. Must be >= 16.
 * @param vlen   Number of vector sources.
 * @param gftbls Pointer to 32*vlen byte array of pre-calculated constants based
 *               on the array of input coefficients. Only elements 32*CONST*j + 1 
 *               of this array are used, where j = (0, 1, 2...) and CONST is the
 *               number of elements in the array of input coefficients. The 
 *               elements used correspond to the original input coefficients.		
 * @param src    Pointers to source input.
 * @param dest   Pointer to destination data array.
 * @returns none
 */

void gf_vect_dot_prod(int len, int vlen, unsigned char *gftbls,
                        unsigned char *src, unsigned char *dest);

/**********************************************************************
 * The remaining are lib support functions used in GF(2^8) operations.
 */

/**
 * @brief Single element GF(2^8) multiply.
 *
 * @param a  Multiplicand a
 * @param b  Multiplicand b
 * @returns  Product of a and b in GF(2^8)
 */

unsigned char gf_mul(unsigned char a, unsigned char b);

/**
 * @brief Single element GF(2^8) inverse.
 *
 * @param a  Input element
 * @returns  Field element b such that a x b = {1}
 */

unsigned char gf_inv(unsigned char a);

/**
 * @brief Generate a matrix of coefficients to be used for encoding.
 *
 * Vandermonde matrix example of encoding coefficients where high portion of
 * matrix is identity matrix I and lower portion is constructed as 2^{i*(j-k+1)}
 * i:{0,k-1} j:{k,m-1}. Commonly used method for choosing coefficients in
 * erasure encoding but does not guarantee invertable for every sub matrix.  For
 * large k it is possible to find cases where the decode matrix chosen from
 * sources and parity not in erasure are not invertable. Users may want to
 * adjust for k > 5.
 *
 * @param a  [mxk] array to hold coefficients
 * @param m  number of rows in matrix corresponding to srcs + parity.
 * @param k  number of columns in matrix corresponding to srcs.
 * @returns  none
 */

void gf_gen_rs_matrix(unsigned char *a, int m, int k);

/**
 * @brief Generate a Cauchy matrix of coefficients to be used for encoding.
 *
 * Cauchy matrix example of encoding coefficients where high portion of matrix
 * is identity matrix I and lower portion is constructed as 1/(i + j) | i != j,
 * i:{0,k-1} j:{k,m-1}.  Any sub-matrix of a Cauchy matrix should be invertable.
 *
 * @param a  [mxk] array to hold coefficients
 * @param m  number of rows in matrix corresponding to srcs + parity.
 * @param k  number of columns in matrix corresponding to srcs.
 * @returns  none
 */

void gf_gen_cauchy1_matrix(unsigned char *a, int m, int k);

/**
 * @brief Invert a matrix in GF(2^8)
 *
 * @param in  input matrix
 * @param out output matrix such that [in] x [out] = [I] - identity matrix
 * @param n   size of matrix [nxn]
 * @returns 0 successful, other fail on singular input matrix
 */

int gf_invert_matrix(unsigned char *in, unsigned char *out, const int n);

/**
 * @brief Generate decode matrix from encode matrix
 *
 * @param encode_matrix  input matrix to generate decode matrix
 * @param decode_matrix  output matrix to decode orignal source
 * @param decode_index  ouptut mapping of decode matrix rows to encoded rows
 * @param src_err_list input error list of the form [0, 0, 1, 0, ...] for valid encoded blocks
 * @param src_in_err input error list of the form [2, 5, 6] for invalid encoded blocks
 * @param nerrs number of encoded rows with errors
 * @param nsrcerrs number of source errors less than k
 * @param k number of source rows needed
 * @param m number of encoded rows
 * @returns 0 successful, other fail on singular input matrix
 */
int gf_gen_decode_matrix(unsigned char *encode_matrix,
				unsigned char *decode_matrix,
				unsigned int *decode_index,
				unsigned char *src_err_list,
				unsigned char *src_in_err,
				int nerrs, int nsrcerrs, int k, int m);

/*************************************************************/

#ifdef __cplusplus
}
#endif

#endif //_ERASURE_CODE_H_
