class CreateExtension < ActiveRecord::Migration[6.1]
  def change
    enable_extension 'uuid-ossp'
    create_table :users, id: :uuid do |t|
      t.string :user_id, null: false
      t.string :email, null: false
      t.string :password, null:false
      t.string :username, null: false
      t.boolean :is_valid, default: true

      t.timestamps
  end

end
