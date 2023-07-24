import { Pagination } from 'react-bootstrap';
import { LinkContainer } from 'react-router-bootstrap';

type Props = {
  pages: number;
  page: number;
  isAdmin?: boolean;
  keyword: string;
};

const Paginate = ({ pages, page, isAdmin = false, keyword = '' }: Props) => {
  const pageNumbers = Array.from({ length: pages }, (_, i) => i + 1);
  return (
    <>
      {pages > 1 && (
        <Pagination>
           {pageNumbers.map((pageNumber) => (
            <LinkContainer
              key={pageNumber}
              to={
                !isAdmin
                  ? keyword
                    ? `/search/${keyword}/page/${pageNumber}`
                    : `/page/${pageNumber}`
                  : `/dashboard/product-list/${pageNumber}`
              }
            >
              <Pagination.Item active={pageNumber === page}>{pageNumber}</Pagination.Item>
            </LinkContainer>
          ))}
        </Pagination>
      )}
    </>
  );
};

export default Paginate;
