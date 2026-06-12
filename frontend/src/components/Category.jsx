import './Category.css';

function Category( {label} ) {

  return (
    <>
      <div className="category-oval">
        <span className="label">{label}</span>
      </div>
    </>
  )
}

export default Category
