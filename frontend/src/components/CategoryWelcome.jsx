import './CategoryWelcome.css';

function CategoryWelcome({ icon, label }) {
  return (
    <>
      <div className="category-oval">
        <img src={icon} className="icon" alt="" />
        <span className="label">{label}</span>
      </div>
    </>
  )
}

export default CategoryWelcome